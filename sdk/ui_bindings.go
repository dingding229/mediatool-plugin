package sdk

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
)

func (r *Runtime) bindUI(ctx context.Context, add func(contracts.Descriptor, func(context.Context, json.RawMessage) (any, error)) error) error {
	bind := func(id string, call func(context.Context, json.RawMessage) (any, error)) error {
		return add(contracts.Descriptor{Kind: "ui", ID: id, RequiredScopes: []string{"ui.json"}}, call)
	}
	if p, ok := r.plugin.(FormProvider); ok && r.manifest.Extensions.Forms {
		if e := bind("form", func(c context.Context, _ json.RawMessage) (any, error) { return p.GetForm(c) }); e != nil {
			return e
		}
	}
	if p, ok := r.plugin.(PageProvider); ok {
		for _, key := range r.manifest.Extensions.Pages {
			if e := bind("page."+key, func(c context.Context, _ json.RawMessage) (any, error) { return p.GetPage(c, PageRequest{Key: key}) }); e != nil {
				return e
			}
		}
	}
	if p, ok := r.plugin.(DashboardProvider); ok {
		for _, key := range r.manifest.Extensions.Dashboards {
			if e := bind("dashboard."+key, func(c context.Context, _ json.RawMessage) (any, error) {
				return p.GetDashboard(c, DashboardRequest{Key: key})
			}); e != nil {
				return e
			}
		}
	}
	if p, ok := r.plugin.(DashboardMetaProvider); ok {
		if e := bind("meta", func(c context.Context, _ json.RawMessage) (any, error) { return p.GetDashboardMeta(c) }); e != nil {
			return e
		}
	}
	return nil
}
