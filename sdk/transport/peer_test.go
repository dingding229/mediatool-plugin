package transport

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func pair(t *testing.T, opts Options, register func(*Peer, *Peer)) (*Peer, *Peer) {
	t.Helper()
	a, b := net.Pipe()
	opts.Prefix = "h"
	h, e := New(a, opts)
	if e != nil {
		t.Fatal(e)
	}
	opts.Prefix = "p"
	p, e := New(b, opts)
	if e != nil {
		t.Fatal(e)
	}
	register(h, p)
	ctx, cancel := context.WithCancel(context.Background())
	if e = h.Start(ctx); e != nil {
		t.Fatal(e)
	}
	if e = p.Start(ctx); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() {
		cancel()
		if e := h.Close(); e != nil {
			t.Error(e)
		}
		if e := p.Close(); e != nil {
			t.Error(e)
		}
	})
	return h, p
}
func handle(t *testing.T, p *Peer, m string, h Handler) {
	t.Helper()
	if e := p.Handle(m, h); e != nil {
		t.Fatal(e)
	}
}
func code(t *testing.T, e error, c contracts.Code) {
	t.Helper()
	var rpc *contracts.Error
	if !errors.As(e, &rpc) || rpc.Code != c {
		t.Fatalf("got %v want %d", e, c)
	}
}
func TestBidirectionalCallbackAndPing(t *testing.T) {
	h, p := pair(t, Options{}, func(h, p *Peer) {
		handle(t, h, "host.call", func(ctx context.Context, b json.RawMessage) (any, *contracts.Error) {
			return struct {
				Value int `json:"value"`
			}{42}, nil
		})
		handle(t, p, "plugin.init", func(ctx context.Context, b json.RawMessage) (any, *contracts.Error) {
			var out struct {
				Value int `json:"value"`
			}
			if e := p.Call(ctx, "host.call", struct{}{}, &out); e != nil {
				return nil, contracts.NewError(contracts.InternalError)
			}
			return out, nil
		})
	})
	var out struct {
		Value int `json:"value"`
	}
	if e := h.Call(context.Background(), "plugin.init", struct{}{}, &out); e != nil || out.Value != 42 {
		t.Fatalf("%v %#v", e, out)
	}
	var pong struct {
		Nonce string `json:"nonce"`
	}
	if e := p.Call(context.Background(), "rpc.ping", map[string]string{"nonce": "test"}, &pong); e != nil || pong.Nonce != "test" {
		t.Fatalf("ping %v %#v", e, pong)
	}
}
func TestCancelBusyPanicAndErrorRedaction(t *testing.T) {
	entered := make(chan struct{})
	cancelled := make(chan struct{})
	h, _ := pair(t, Options{MaxHandlers: 1}, func(h, p *Peer) {
		handle(t, p, "wait", func(ctx context.Context, b json.RawMessage) (any, *contracts.Error) {
			close(entered)
			<-ctx.Done()
			close(cancelled)
			return nil, nil
		})
		handle(t, p, "panic", func(ctx context.Context, b json.RawMessage) (any, *contracts.Error) {
			panic("secret-url?token=private")
		})
		handle(t, p, "unsafe", func(ctx context.Context, b json.RawMessage) (any, *contracts.Error) {
			return nil, &contracts.Error{Code: contracts.NotFound, Message: "secret-url?token=private", Data: contracts.ErrorData{Reason: "private"}}
		})
	})
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() { result <- h.Call(ctx, "wait", nil, nil) }()
	<-entered
	code(t, h.Call(context.Background(), "wait", nil, nil), contracts.Busy)
	var pong string
	if e := h.Call(context.Background(), "rpc.ping", "alive", &pong); e != nil || pong != "alive" {
		t.Fatal(e)
	}
	cancel()
	if e := <-result; !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	select {
	case <-cancelled:
	case <-time.After(time.Second):
		t.Fatal("remote handler not cancelled")
	}
	// Wait until the callback slot has been released, without assuming goroutine timing.
	for deadline := time.Now().Add(time.Second); ; {
		e := h.Call(context.Background(), "unsafe", nil, nil)
		var x *contracts.Error
		if errors.As(e, &x) && x.Code == contracts.NotFound {
			raw, _ := json.Marshal(x)
			if strings.Contains(string(raw), "private") {
				t.Fatal("error leaked")
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal(e)
		}
		time.Sleep(time.Millisecond)
	}
	code(t, h.Call(context.Background(), "panic", nil, nil), contracts.InternalError)
}
func TestTimeoutPendingBoundAndCloseUnblock(t *testing.T) {
	entered := make(chan struct{})
	h, _ := pair(t, Options{MaxPending: 1, CallTimeout: 50 * time.Millisecond}, func(h, p *Peer) {
		handle(t, p, "wait", func(ctx context.Context, b json.RawMessage) (any, *contracts.Error) {
			close(entered)
			<-ctx.Done()
			return nil, nil
		})
	})
	done := make(chan error, 1)
	go func() { done <- h.Call(context.Background(), "wait", nil, nil) }()
	<-entered
	code(t, h.Call(context.Background(), "rpc.ping", nil, nil), contracts.Busy)
	if e := <-done; !errors.Is(e, context.DeadlineExceeded) {
		t.Fatal(e)
	}
}
func TestParallelCalls(t *testing.T) {
	h, _ := pair(t, Options{MaxPending: 32}, func(h, p *Peer) {
		handle(t, p, "echo", func(ctx context.Context, b json.RawMessage) (any, *contracts.Error) { return b, nil })
	})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Go(func() {
			for j := 0; j < 30; j++ {
				var out int
				if e := h.Call(context.Background(), "echo", j, &out); e != nil || out != j {
					t.Errorf("%v %d", e, out)
					return
				}
			}
		})
	}
	wg.Wait()
}
func rawFrame(b string) []byte {
	out := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(out, uint32(len(b)))
	copy(out[4:], b)
	return out
}
func TestFrameValidation(t *testing.T) {
	good := ` {"jsonrpc":"2.0","id":"h:1","method":"a","params":{}} `
	tests := []struct {
		name, body string
		ok         bool
	}{
		{"valid", good, true}, {"batch", `[]`, false}, {"duplicate", `{"jsonrpc":"2.0","method":"a","params":{"x":1,"x":2}}`, false},
		{"trailing", good + `{}`, false}, {"unknown", `{"jsonrpc":"2.0","method":"a","extra":1}`, false}, {"both", `{"jsonrpc":"2.0","id":"p:1","result":{},"error":{"code":-32603}}`, false},
		{"nullID", `{"jsonrpc":"2.0","id":null,"result":{}}`, false}, {"numericID", `{"jsonrpc":"2.0","id":1,"method":"a"}`, false}, {"depth", `{"jsonrpc":"2.0","method":"a","params":` + strings.Repeat("[", 65) + `0` + strings.Repeat("]", 65) + `}`, false},
		{"nan", `{"jsonrpc":"2.0","method":"a","params":NaN}`, false}, {"invalidUTF8", "{\xff}", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := readFrame(bytes.NewReader(rawFrame(tt.body)), 1<<20)
			if (e == nil) != tt.ok {
				t.Fatalf("valid=%v error=%v", tt.ok, e)
			}
		})
	}
	oversized := make([]byte, 4)
	binary.BigEndian.PutUint32(oversized, 10000000)
	_, e := readFrame(bytes.NewReader(oversized), 128)
	code(t, e, contracts.PayloadTooLarge)
	_, e = readFrame(bytes.NewReader([]byte{0, 0}), 128)
	if !errors.Is(e, io.ErrUnexpectedEOF) {
		t.Fatal(e)
	}
	_, e = readFrame(bytes.NewReader(append([]byte{0, 0, 0, 2}, '{')), 128)
	if !errors.Is(e, io.ErrUnexpectedEOF) {
		t.Fatal(e)
	}
	joined := append(rawFrame(good), rawFrame(good)...)
	r := bytes.NewReader(joined)
	for range 2 {
		if _, e = readFrame(r, 1<<20); e != nil {
			t.Fatal(e)
		}
	}
}
func TestPartialIOAndMalformedResponseClose(t *testing.T) {
	a, b := net.Pipe()
	h, e := New(a, Options{Prefix: "h"})
	if e != nil {
		t.Fatal(e)
	}
	if e = h.Start(context.Background()); e != nil {
		t.Fatal(e)
	}
	defer h.Close()
	defer b.Close()
	body := rawFrame(`{"jsonrpc":"2.0","id":"p:1","method":"rpc.ping","params":{"nonce":1}}`)
	finished := make(chan error, 1)
	go func() {
		for _, x := range body {
			if _, e := b.Write([]byte{x}); e != nil {
				finished <- e
				return
			}
		}
		finished <- nil
	}()
	m, e := readFrame(b, 1<<20)
	if e != nil || m.ID != "p:1" {
		t.Fatal(e)
	}
	if e = <-finished; e != nil {
		t.Fatal(e)
	}
	go func() { _, _ = b.Write(rawFrame(`{"jsonrpc":"2.0","id":"wrong","result":{}}`)) }()
	select {
	case <-h.Done():
	case <-time.After(time.Second):
		t.Fatal("invalid response not closed")
	}
}
func TestOutputLimitAndNotifications(t *testing.T) {
	var notified atomic.Int32
	h, _ := pair(t, Options{MaxFrame: 256}, func(h, p *Peer) {
		handle(t, p, "huge", func(context.Context, json.RawMessage) (any, *contracts.Error) { return strings.Repeat("a", 1024), nil })
		handle(t, p, "note", func(context.Context, json.RawMessage) (any, *contracts.Error) { notified.Add(1); return nil, nil })
	})
	code(t, h.Call(context.Background(), "huge", nil, nil), contracts.PayloadTooLarge)
	if e := h.Notify("note", nil); e != nil {
		t.Fatal(e)
	}
	for deadline := time.Now().Add(time.Second); notified.Load() == 0; {
		if time.Now().After(deadline) {
			t.Fatal("notification not delivered")
		}
		time.Sleep(time.Millisecond)
	}
	code(t, h.Call(context.Background(), "x", strings.Repeat("a", 1024), nil), contracts.PayloadTooLarge)
	code(t, h.Call(context.Background(), "unknown", nil, nil), contracts.MethodNotFound)
}
func TestOptionsAndLifecycle(t *testing.T) {
	for _, o := range []Options{{Prefix: "bad"}, {Prefix: "h", MaxFrame: 5 << 20}, {Prefix: "h", MaxHandlers: 5}, {Prefix: "h", Queue: 100}, {Prefix: "h", CallTimeout: time.Hour}} {
		a, b := net.Pipe()
		if _, e := New(a, o); e == nil {
			t.Fatal(o)
		}
		a.Close()
		b.Close()
	}
	h, _ := pair(t, Options{}, func(h, p *Peer) {})
	code(t, h.Start(context.Background()), contracts.Conflict)
	code(t, h.Handle("late", func(context.Context, json.RawMessage) (any, *contracts.Error) { return nil, nil }), contracts.Conflict)
}

type panicMarshaler struct{}

func (panicMarshaler) MarshalJSON() ([]byte, error) { panic("private renderer bug") }
func TestMarshalerPanicAndCancelledCallerBeforeSend(t *testing.T) {
	h, _ := pair(t, Options{}, func(h, p *Peer) {
		handle(t, p, "marshalpanic", func(context.Context, json.RawMessage) (any, *contracts.Error) { return panicMarshaler{}, nil })
	})
	code(t, h.Call(context.Background(), "marshalpanic", nil, nil), contracts.InternalError)
	code(t, h.Call(context.Background(), "rpc.ping", panicMarshaler{}, nil), contracts.InvalidParams)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := h.Call(ctx, "rpc.ping", nil, nil); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
func TestQueuedCancellationCannotOvertakeItsRequest(t *testing.T) {
	a, b := net.Pipe()
	h, e := New(a, Options{Prefix: "h"})
	if e != nil {
		t.Fatal(e)
	}
	if e = h.Start(context.Background()); e != nil {
		t.Fatal(e)
	}
	defer h.Close()
	defer b.Close()
	// The first queued packet blocks in Write until the remote starts reading.
	if e = h.Notify("note", nil); e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- h.Call(ctx, "wait", nil, nil) }()
	for deadline := time.Now().Add(time.Second); ; {
		h.mu.Lock()
		n := len(h.pending)
		h.mu.Unlock()
		if n == 1 && len(h.normal) > 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("request not queued")
		}
		time.Sleep(time.Millisecond)
	}
	cancel()
	if e = <-done; !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	got := make(chan []string, 1)
	go func() {
		names := []string{}
		for range 3 {
			m, e := readFrame(b, 1<<20)
			if e != nil {
				got <- nil
				return
			}
			names = append(names, m.Method)
		}
		got <- names
	}()
	select {
	case names := <-got:
		if strings.Join(names, ",") != "note,wait,$/cancelRequest" {
			t.Fatal(names)
		}
	case <-time.After(time.Second):
		t.Fatal("queued cancellation blocked")
	}
}
func TestByteAndQueueBackpressure(t *testing.T) {
	a, b := net.Pipe()
	h, e := New(a, Options{Prefix: "h", Queue: 1, MaxFrame: 256, MaxQueuedBytes: 256})
	if e != nil {
		t.Fatal(e)
	}
	if e = h.Start(context.Background()); e != nil {
		t.Fatal(e)
	}
	defer h.Close()
	defer b.Close()
	if e = h.Notify("note", strings.Repeat("x", 120)); e != nil {
		t.Fatal(e)
	}
	code(t, h.Notify("note", strings.Repeat("x", 120)), contracts.Busy)
	if h.queued.Load() > 256 {
		t.Fatal("byte budget exceeded")
	}
	if e = h.Close(); e != nil {
		t.Fatal(e)
	}
	select {
	case <-h.Done():
	case <-time.After(time.Second):
		t.Fatal("blocked writer leaked")
	}
}
