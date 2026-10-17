package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"io"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Handler func(context.Context, json.RawMessage) (any, *contracts.Error)
type Options struct {
	Prefix                                                   string
	MaxFrame, Queue, MaxHandlers, MaxPending, MaxQueuedBytes int
	CallTimeout, HandlerTimeout                              time.Duration
}
type queuedFrame struct {
	data []byte
	sent *atomic.Bool
}
type Peer struct {
	ready               chan struct{}
	conn                io.ReadWriteCloser
	opts                Options
	remote              string
	mu                  sync.Mutex
	handlers            map[string]Handler
	pending             map[string]chan message
	active              map[string]context.CancelFunc
	started             bool
	err                 error
	root                context.Context
	cancel              context.CancelFunc
	done                chan struct{}
	high, normal        chan queuedFrame
	slots, pendingSlots chan struct{}
	sequence            atomic.Uint64
	queued              atomic.Int64
	once                sync.Once
}

func New(conn io.ReadWriteCloser, o Options) (*Peer, error) {
	if conn == nil || (o.Prefix != "h" && o.Prefix != "p") {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	if o.MaxFrame == 0 {
		o.MaxFrame = 1 << 20
	}
	if o.Queue == 0 {
		o.Queue = 32
	}
	if o.MaxHandlers == 0 {
		o.MaxHandlers = 4
	}
	if o.MaxPending == 0 {
		o.MaxPending = 8
	}
	if o.MaxQueuedBytes == 0 {
		o.MaxQueuedBytes = 4 << 20
	}
	if o.CallTimeout == 0 {
		o.CallTimeout = 30 * time.Second
	}
	if o.HandlerTimeout == 0 {
		o.HandlerTimeout = 5 * time.Minute
	}
	if o.MaxFrame < 128 || o.MaxFrame > contracts.MaxJSON || o.Queue < 1 || o.Queue > 32 || o.MaxHandlers < 1 || o.MaxHandlers > 4 || o.MaxPending < 1 || o.MaxPending > 32 || o.MaxQueuedBytes < o.MaxFrame || o.MaxQueuedBytes > 4<<20 || o.CallTimeout <= 0 || o.CallTimeout > 5*time.Minute || o.HandlerTimeout <= 0 || o.HandlerTimeout > 5*time.Minute {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	remote := "h"
	if o.Prefix == "h" {
		remote = "p"
	}
	return &Peer{conn: conn, opts: o, remote: remote, handlers: map[string]Handler{}, pending: map[string]chan message{}, active: map[string]context.CancelFunc{}, done: make(chan struct{}), ready: make(chan struct{}), high: make(chan queuedFrame, 8), normal: make(chan queuedFrame, o.Queue), slots: make(chan struct{}, o.MaxHandlers), pendingSlots: make(chan struct{}, o.MaxPending)}, nil
}

// Handle must be called before Start; callbacks never execute on the IO reader.
func (p *Peer) Handle(method string, h Handler) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.started || method == "" || len(method) > 128 || h == nil || method == "rpc.ping" || method == "$/cancelRequest" || p.handlers[method] != nil {
		return contracts.NewError(contracts.Conflict)
	}
	p.handlers[method] = h
	return nil
}
func (p *Peer) Start(ctx context.Context) error {
	if ctx == nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	p.mu.Lock()
	if p.started || p.err != nil {
		p.mu.Unlock()
		return contracts.NewError(contracts.Conflict)
	}
	p.started = true
	p.root, p.cancel = context.WithCancel(ctx)
	root := p.root
	close(p.ready)
	p.mu.Unlock()
	go p.readLoop()
	go p.writeLoop()
	go func() {
		select {
		case <-root.Done():
			p.shutdown(contracts.NewError(contracts.Cancelled))
		case <-p.done:
		}
	}()
	return nil
}
func (p *Peer) Ready() <-chan struct{} { return p.ready }
func (p *Peer) Done() <-chan struct{}  { return p.done }
func (p *Peer) Err() error             { p.mu.Lock(); defer p.mu.Unlock(); return p.err }
func (p *Peer) Close() error           { p.shutdown(contracts.NewError(contracts.Unavailable)); return nil }
func (p *Peer) shutdown(err error) {
	p.once.Do(func() {
		p.mu.Lock()
		p.err = err
		if p.cancel != nil {
			p.cancel()
		}
		for _, cancel := range p.active {
			cancel()
		}
		close(p.done)
		p.mu.Unlock()
		// Closing IO is necessary to unblock both a partial frame and a blocked writer.
		_ = p.conn.Close() // Teardown only: the primary transport error is already retained.
	})
}
func (p *Peer) enqueue(m message, priority bool, sent ...*atomic.Bool) error {
	p.mu.Lock()
	ready := p.started && p.err == nil
	p.mu.Unlock()
	if !ready {
		return contracts.NewError(contracts.Unavailable)
	}
	b, e := encode(m, p.opts.MaxFrame)
	if e != nil {
		return e
	}
	n := int64(len(b) - 4)
	for {
		old := p.queued.Load()
		if old+n > int64(p.opts.MaxQueuedBytes) {
			return contracts.NewError(contracts.Busy)
		}
		if p.queued.CompareAndSwap(old, old+n) {
			break
		}
	}
	q := p.normal
	if priority {
		q = p.high
	}
	packet := queuedFrame{data: b}
	if len(sent) > 0 {
		packet.sent = sent[0]
	}
	select {
	case q <- packet:
		return nil
	case <-p.done:
		p.queued.Add(-n)
		return contracts.NewError(contracts.Unavailable)
	default:
		p.queued.Add(-n)
		return contracts.NewError(contracts.Busy)
	}
}
func (p *Peer) Call(ctx context.Context, method string, args, out any) error {
	if ctx == nil || method == "" || len(method) > 128 {
		return contracts.NewError(contracts.InvalidParams)
	}
	ctx, cancel := context.WithTimeout(ctx, p.opts.CallTimeout)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case p.pendingSlots <- struct{}{}:
		defer func() { <-p.pendingSlots }()
	default:
		return contracts.NewError(contracts.Busy)
	}
	raw, e := safeMarshal(args)
	if e != nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	if contracts.ValidateJSON(raw) != nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	id := p.opts.Prefix + ":" + strconv.FormatUint(p.sequence.Add(1), 10)
	ch := make(chan message, 1)
	p.mu.Lock()
	if p.err != nil || !p.started {
		p.mu.Unlock()
		return contracts.NewError(contracts.Unavailable)
	}
	p.pending[id] = ch
	p.mu.Unlock()
	defer func() { p.mu.Lock(); delete(p.pending, id); p.mu.Unlock() }()
	var sent atomic.Bool
	if e = p.enqueue(message{Version: "2.0", ID: id, Method: method, Params: raw}, false, &sent); e != nil {
		return e
	}
	select {
	case m := <-ch:
		if m.Error != nil {
			return contracts.SafeError(m.Error)
		}
		if out != nil {
			return contracts.Decode(m.Result, out)
		}
		return nil
	case <-ctx.Done():
		if e = p.notify("$/cancelRequest", struct {
			ID string `json:"id"`
		}{id}, sent.Load()); e != nil {
			p.shutdown(contracts.NewError(contracts.Cancelled))
		}
		return ctx.Err()
	case <-p.done:
		return p.Err()
	}
}
func (p *Peer) Notify(method string, args any) error { return p.notify(method, args, false) }
func (p *Peer) notify(method string, args any, priority bool) error {
	if method == "" || len(method) > 128 {
		return contracts.NewError(contracts.InvalidParams)
	}
	b, e := safeMarshal(args)
	if e != nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	return p.enqueue(message{Version: "2.0", Method: method, Params: b}, priority)
}
func (p *Peer) writeLoop() {
	for {
		var packet queuedFrame
		select {
		case <-p.done:
			return
		case packet = <-p.high:
		default:
			select {
			case <-p.done:
				return
			case packet = <-p.high:
			case packet = <-p.normal:
			}
		}
		if e := writeFull(p.conn, packet.data); e != nil {
			p.shutdown(contracts.NewError(contracts.Unavailable))
			return
		}
		if packet.sent != nil {
			packet.sent.Store(true)
		}
		p.queued.Add(-int64(len(packet.data) - 4))
	}
}
func (p *Peer) readLoop() {
	for {
		m, e := readFrame(p.conn, p.opts.MaxFrame)
		if e != nil {
			p.shutdown(contracts.NewError(contracts.Unavailable))
			return
		}
		if m.Method == "" {
			if !validID(m.ID, p.opts.Prefix) {
				p.shutdown(contracts.NewError(contracts.InvalidRequest))
				return
			}
			p.mu.Lock()
			ch := p.pending[m.ID]
			if ch != nil {
				select {
				case ch <- m:
				default:
				}
			}
			p.mu.Unlock()
			continue
		}
		if m.ID != "" && !validID(m.ID, p.remote) {
			p.shutdown(contracts.NewError(contracts.InvalidRequest))
			return
		}
		if m.Method == "$/cancelRequest" {
			var c struct {
				ID string `json:"id"`
			}
			if m.ID != "" || contracts.Decode(m.Params, &c) != nil || !validID(c.ID, p.remote) {
				p.shutdown(contracts.NewError(contracts.InvalidRequest))
				return
			}
			p.mu.Lock()
			cancel := p.active[c.ID]
			p.mu.Unlock()
			if cancel != nil {
				cancel()
			}
			continue
		}
		if m.Method == "rpc.ping" {
			p.reply(m, m.Params, nil)
			continue
		}
		p.dispatch(m)
	}
}
func (p *Peer) reply(m message, result any, e *contracts.Error) {
	if m.ID == "" {
		return
	}
	if e != nil {
		e = contracts.SafeError(e)
	}
	r := message{Version: "2.0", ID: m.ID, Error: e}
	if e == nil {
		var err error
		r.Result, err = safeMarshal(result)
		if err != nil {
			r.Error = contracts.NewError(contracts.InternalError)
			r.Result = nil
		}
	}
	if err := p.enqueue(r, true); err != nil {
		var pe *contracts.Error
		if errors.As(err, &pe) && pe.Code == contracts.PayloadTooLarge {
			r.Result = nil
			r.Error = contracts.NewError(contracts.PayloadTooLarge)
			if p.enqueue(r, true) == nil {
				return
			}
		}
		p.shutdown(contracts.NewError(contracts.Unavailable))
	}
}
func (p *Peer) dispatch(m message) {
	p.mu.Lock()
	h := p.handlers[m.Method]
	_, duplicate := p.active[m.ID]
	p.mu.Unlock()
	if duplicate && m.ID != "" {
		p.shutdown(contracts.NewError(contracts.InvalidRequest))
		return
	}
	if h == nil {
		p.reply(m, nil, contracts.NewError(contracts.MethodNotFound))
		return
	}
	select {
	case p.slots <- struct{}{}:
	default:
		p.reply(m, nil, contracts.NewError(contracts.Busy))
		return
	}
	ctx, cancel := context.WithTimeout(p.root, p.opts.HandlerTimeout)
	p.mu.Lock()
	if m.ID != "" {
		p.active[m.ID] = cancel
	}
	p.mu.Unlock()
	go func() {
		cleanup := sync.OnceFunc(func() { cancel(); p.mu.Lock(); delete(p.active, m.ID); p.mu.Unlock(); <-p.slots })
		defer cleanup()
		var result any
		var e *contracts.Error
		func() {
			defer func() {
				if recover() != nil {
					e = contracts.NewError(contracts.InternalError)
				}
			}()
			result, e = h(ctx, m.Params)
		}()
		if ctx.Err() != nil {
			e = contracts.NewError(contracts.Cancelled)
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				e = contracts.NewError(contracts.DeadlineExceeded)
			}
		}
		cleanup()
		p.reply(m, result, e)
	}()
}
