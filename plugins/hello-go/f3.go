package hello

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/host/cache"
	"github.com/dingding229/mediatool-plugin/sdk/host/downloads"
	"github.com/dingding229/mediatool-plugin/sdk/host/integrations"
	"github.com/dingding229/mediatool-plugin/sdk/host/jobs"
	"github.com/dingding229/mediatool-plugin/sdk/host/library"
	"github.com/dingding229/mediatool-plugin/sdk/host/maintenance"
	"github.com/dingding229/mediatool-plugin/sdk/host/mediaserver"
	"github.com/dingding229/mediatool-plugin/sdk/host/network"
	"github.com/dingding229/mediatool-plugin/sdk/host/notifications"
	"github.com/dingding229/mediatool-plugin/sdk/host/settings"
	"github.com/dingding229/mediatool-plugin/sdk/host/system"
)

// Manual read-only acceptance: outputs contain only booleans, never data from
// the host. The denied write MUST fail before settings validation/effect.
func (p *Plugin) f3ReadCheck(ctx context.Context, q sdk.ActionContext) (sdk.ActionContext, bool, error) {
	page := b.PageRequest{Page: 1, Size: 1}
	checks := map[string]bool{}
	_, e := downloads.New(p.host).List(ctx, page)
	checks["F3-A.downloads"] = e == nil
	_, e = library.New(p.host).Assets(ctx, page)
	checks["F3-B.library"] = e == nil
	_, e = mediaserver.New(p.host).Providers(ctx, b.Empty{})
	checks["F3-C.mediaserver"] = e == nil
	_, e = cache.New(p.host).Status(ctx, b.NamespaceQuery{Namespace: "tmdb", PageRequest: page})
	checks["F3-D.cache"] = e == nil
	_, e = integrations.New(p.host).List(ctx, page)
	checks["F3-D.integrations"] = e == nil
	_, e = jobs.New(p.host).List(ctx, b.JobQuery{PageRequest: page})
	checks["F3-D.jobs"] = e == nil
	_, e = maintenance.New(p.host).Health(ctx, b.Empty{})
	checks["F3-D.maintenance"] = e == nil
	_, e = notifications.New(p.host).Deliveries(ctx, page)
	checks["F3-D.notifications"] = e == nil
	_, e = network.New(p.host).ProxyStatus(ctx, b.Empty{})
	checks["F3-D.networkproxy"] = e == nil
	_, e = settings.New(p.host).Get(ctx, b.SettingKeys{Keys: []string{"searchCacheSeconds"}})
	checks["F3-D.settings"] = e == nil
	_, e = system.New(p.host).Info(ctx, b.Empty{})
	checks["F3-D.system"] = e == nil
	_, e = settings.New(p.host).Patch(ctx, b.SettingPatch{Revision: 1, Values: map[string]json.RawMessage{"searchCacheSeconds": json.RawMessage(`900`)}}, b.Operation{Key: "hello-denied-write", Confirm: true})
	var rpc *contracts.Error
	checks["denied.settings.write"] = errors.As(e, &rpc) && rpc.Code == contracts.PermissionDenied
	passed := true
	for _, ok := range checks {
		passed = passed && ok
	}
	result := struct {
		Passed bool            `json:"passed"`
		Checks map[string]bool `json:"checks"`
	}{passed, checks}
	var previous any
	revision, e := p.host.GetData(ctx, "f3-read-check", &previous)
	if e != nil && (!errors.As(e, &rpc) || rpc.Code != contracts.NotFound) {
		return q, false, e
	}
	if _, e = p.host.SaveData(ctx, "f3-read-check", result, revision); e != nil {
		return q, false, e
	}
	if q.Data == nil {
		q.Data = map[string]json.RawMessage{}
	}
	raw, _ := json.Marshal(result)
	q.Data["f3-read-check"] = raw
	return q, true, nil
}
