package sdk

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	"github.com/dingding229/mediatool-plugin/sdk/schema"
	"github.com/dingding229/mediatool-plugin/sdk/transport"
	"reflect"
	"slices"
	"strings"
	"time"
)

type Identity struct{ InstanceID, DefinitionFingerprint, CatalogBuild string }
type Runtime struct {
	peer        *transport.Peer
	manifest    contracts.Manifest
	factory     Factory
	identity    Identity
	lane        chan struct{}
	initialized bool
	started     bool
	failed      bool
	stopped     bool
	init        contracts.Initialize
	plugin      Plugin
	host        *host.Client
	events      *events.Registry
	bindings    map[string]bound
}

func NewRuntime(peer *transport.Peer, m contracts.Manifest, f Factory, id Identity) (*Runtime, error) {
	if peer == nil || f == nil || id.InstanceID == "" || id.DefinitionFingerprint == "" || id.CatalogBuild == "" {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	raw, e := json.Marshal(m)
	if e != nil {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	m, e = contracts.ParseManifest(raw)
	if e != nil {
		return nil, e
	}
	r := &Runtime{peer: peer, manifest: m, factory: f, identity: id, host: host.New(peer), events: events.New(), lane: make(chan struct{}, 1)}
	for _, name := range r.methods() {
		method := name
		if e = peer.Handle(method, func(ctx context.Context, b json.RawMessage) (any, *contracts.Error) { return r.handle(ctx, method, b) }); e != nil {
			return nil, e
		}
	}
	return r, nil
}
func (r *Runtime) methods() []string {
	return []string{"plugin.describe", "plugin.invoke", "plugin.initialize", "plugin.init", "plugin.get_state", "plugin.get_name", "plugin.get_form", "plugin.get_page", "plugin.get_render_mode", "plugin.get_event_handlers", "plugin.get_api", "plugin.get_service", "plugin.get_command", "plugin.get_dashboard", "plugin.get_dashboard_meta", "plugin.get_module", "plugin.get_auth_providers", "plugin.get_actions", "plugin.stop_service", "plugin.close"}
}
func (r *Runtime) On(t events.EventType, h events.Handler) error { return r.events.On(t, h) }
func (r *Runtime) OnChain(t events.ChainEventType, h events.ChainHandler) error {
	return r.events.OnChain(t, h)
}
func (r *Runtime) handle(ctx context.Context, method string, b json.RawMessage) (any, *contracts.Error) {
	budget := 2 * time.Second
	if method == "plugin.initialize" {
		budget = 5 * time.Second
	}
	if method == "plugin.invoke" {
		budget = 30 * time.Second
	}
	if method == "plugin.init" {
		budget = 15 * time.Second
	}
	if method == "plugin.stop_service" || method == "plugin.close" {
		budget = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, budget)
	defer cancel()
	select {
	case r.lane <- struct{}{}:
		defer func() { <-r.lane }()
	case <-ctx.Done():
		return nil, contracts.NewError(contracts.DeadlineExceeded)
	}
	if method == "plugin.initialize" {
		return r.initialize(b)
	}
	if !r.initialized {
		return nil, contracts.NewError(contracts.Unavailable)
	}
	if method == "plugin.stop_service" || method == "plugin.close" {
		if r.stopped {
			return struct{}{}, nil
		}
		r.stopped = true
		if r.plugin != nil {
			if e := r.plugin.StopService(ctx); e != nil {
				return nil, contracts.NewError(contracts.NeedsAction)
			}
		}
		return struct{}{}, nil
	}
	if r.stopped || r.failed {
		return nil, contracts.NewError(contracts.NeedsAction)
	}
	if method == "plugin.init" {
		if r.started {
			return nil, contracts.NewError(contracts.Conflict)
		}
		var config Config
		if contracts.Decode(b, &config) != nil || config.Revision != r.init.ConfigRevision || contracts.ValidateJSON(config.Value) != nil {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		// Set failed BEFORE entering user code, so a recovered panic cannot expose half-init state.
		r.failed = true
		r.plugin = r.factory(Services{r.host, r.events, r.identity.InstanceID})
		if r.plugin == nil || (reflect.ValueOf(r.plugin).Kind() == reflect.Ptr && reflect.ValueOf(r.plugin).IsNil()) {
			return nil, contracts.NewError(contracts.InternalError)
		}
		if e := CheckCapabilities(r.manifest, r.plugin); e != nil {
			return nil, contracts.NewError(contracts.InvalidExtensionResult)
		}
		if e := r.plugin.InitPlugin(ctx, config); e != nil {
			return nil, contracts.NewError(contracts.NeedsAction)
		}
		r.events.Freeze()
		for _, d := range r.events.Descriptors() {
			scope := "events.subscribe." + d.Type
			declared := r.manifest.Extensions.Events
			if d.Kind == "chain" {
				scope = "events.chain." + d.Type
				declared = r.manifest.Extensions.Chains
			}
			if !slices.Contains(declared, d.Type) || !slices.Contains(r.init.Capabilities, scope) {
				return nil, contracts.NewError(contracts.PermissionDenied)
			}
		}
		if e := r.bind(ctx); e != nil {
			return nil, contracts.SafeError(contracts.NewError(contracts.InvalidExtensionResult))
		}
		r.failed = false
		r.started = true
		return struct{}{}, nil
	}
	if !r.started {
		return nil, contracts.NewError(contracts.Unavailable)
	}
	if method == "plugin.invoke" {
		return r.invoke(ctx, b)
	}
	if method == "plugin.describe" {
		return r.descriptors(""), nil
	}
	for methodName, kind := range map[string]string{"plugin.get_api": "api", "plugin.get_service": "service", "plugin.get_command": "command", "plugin.get_module": "module", "plugin.get_actions": "action", "plugin.get_auth_providers": "chain"} {
		if method == methodName {
			return r.descriptors(kind).Descriptors, nil
		}
	}
	var result any
	var err error
	switch method {
	case "plugin.get_state":
		result, err = r.plugin.GetState(ctx)
	case "plugin.get_name":
		result = r.manifest.Name
		if p, ok := r.plugin.(NameProvider); ok {
			result, err = p.GetName(ctx)
		}
	case "plugin.get_render_mode":
		result = struct {
			Mode string `json:"mode"`
		}{r.manifest.UI.Mode}
	case "plugin.get_form":
		if p, ok := r.plugin.(FormProvider); ok && r.manifest.Extensions.Forms {
			result, err = p.GetForm(ctx)
		} else {
			return nil, contracts.NewError(contracts.UnsupportedVersion)
		}
	case "plugin.get_page":
		var q PageRequest
		if contracts.Decode(b, &q) != nil || !slices.Contains(r.manifest.Extensions.Pages, q.Key) {
			return nil, contracts.NewError(contracts.NotFound)
		}
		if p, ok := r.plugin.(PageProvider); ok {
			result, err = p.GetPage(ctx, q)
		} else {
			return nil, contracts.NewError(contracts.UnsupportedVersion)
		}
	case "plugin.get_event_handlers":
		result = r.events.Descriptors()
	case "plugin.get_dashboard":
		var q DashboardRequest
		if contracts.Decode(b, &q) != nil || !slices.Contains(r.manifest.Extensions.Dashboards, q.Key) {
			return nil, contracts.NewError(contracts.NotFound)
		}
		if p, ok := r.plugin.(DashboardProvider); ok {
			result, err = p.GetDashboard(ctx, q)
		} else {
			return nil, contracts.NewError(contracts.UnsupportedVersion)
		}
	case "plugin.get_dashboard_meta":
		if p, ok := r.plugin.(DashboardMetaProvider); ok {
			metas, e := p.GetDashboardMeta(ctx)
			err = e
			if len(metas) > 8 {
				return nil, contracts.NewError(contracts.InvalidExtensionResult)
			}
			seen := map[string]bool{}
			for _, m := range metas {
				if !slices.Contains(r.manifest.Extensions.Dashboards, m.Key) || seen[m.Key] || m.Refresh < 10 || m.Cols < 1 || m.Cols > 12 || len(m.Title) > 120 {
					return nil, contracts.NewError(contracts.InvalidExtensionResult)
				}
				seen[m.Key] = true
			}
			result = metas
		} else {
			return nil, contracts.NewError(contracts.UnsupportedVersion)
		}
	default:
		return nil, contracts.NewError(contracts.MethodNotFound)
	}
	if err != nil {
		return nil, contracts.NewError(contracts.NeedsAction)
	}

	return result, nil
}
func (r *Runtime) initialize(b json.RawMessage) (any, *contracts.Error) {
	var q contracts.Initialize
	if contracts.Decode(b, &q) != nil || len(q.InstanceID) > 128 || len(q.Generation) > 128 || q.Generation == "" || q.ConfigRevision < 0 || len(q.Capabilities) > 128 {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	if q.ContractVersion != ContractVersion || q.CatalogBuild != r.identity.CatalogBuild || q.DefinitionFingerprint != r.identity.DefinitionFingerprint || q.InstanceID != r.identity.InstanceID {
		return nil, contracts.NewError(contracts.UnsupportedVersion)
	}
	seen := map[string]bool{}
	for _, cap := range q.Capabilities {
		if !contracts.KnownScope(cap) || seen[cap] {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		seen[cap] = true
	}
	// Scope declarations and Host grants must BOTH permit a call; never silently grant a declaration.
	for _, p := range r.manifest.Permissions {
		if !seen[p.Scope] {
			return nil, contracts.NewError(contracts.PermissionDenied)
		}
	}
	for cap := range seen {
		declared := false
		for _, p := range r.manifest.Permissions {
			if p.Scope == cap {
				declared = true
				break
			}
		}
		if !declared {
			return nil, contracts.NewError(contracts.PermissionDenied)
		}
	}
	if r.initialized && !reflect.DeepEqual(q, r.init) {
		return nil, contracts.NewError(contracts.Conflict)
	}
	r.init = q
	r.initialized = true
	return contracts.Initialized{ContractVersion: ContractVersion, SchemaFingerprint: schema.Fingerprint, CatalogBuild: r.identity.CatalogBuild, ImplementedMethods: r.methods()}, nil
}

// Serve only starts the SDK transport. G2's supervisor owns process isolation/reaping.
func (r *Runtime) Serve(ctx context.Context) error {
	if e := r.peer.Start(ctx); e != nil {
		return e
	}
	<-r.peer.Done()
	// Cleanup cannot wait forever on a noncooperative callback; G2 reaps the process.
	cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	select {
	case r.lane <- struct{}{}:
	case <-cleanup.Done():
		return contracts.NewError(contracts.NeedsAction)
	}
	result := make(chan error, 1)
	go func() {
		defer func() {
			<-r.lane
			if recover() != nil {
				result <- contracts.NewError(contracts.InternalError)
			}
		}()
		if r.plugin != nil && !r.stopped {
			r.stopped = true
			result <- r.plugin.StopService(cleanup)
		} else {
			result <- nil
		}
	}()
	select {
	case e := <-result:
		if e != nil {
			return contracts.NewError(contracts.NeedsAction)
		}
		return r.peer.Err()
	case <-cleanup.Done():
		return contracts.NewError(contracts.NeedsAction)
	}
}

// CheckCapabilities is useful to validate a factory's optional-interface declarations.
func CheckCapabilities(m contracts.Manifest, p Plugin) error {
	if m.Extensions.Forms {
		if _, ok := p.(FormProvider); !ok {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
	}
	if len(m.Extensions.Pages) > 0 {
		if _, ok := p.(PageProvider); !ok {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
	}
	for _, permission := range m.Permissions {
		if strings.Contains(permission.Scope, "*") {
			return contracts.NewError(contracts.InvalidParams)
		}
	}
	return nil
}
