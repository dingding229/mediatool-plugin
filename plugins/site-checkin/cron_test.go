package checkin

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"strings"
	"testing"
)

func TestCronFormsAndSnapshotFieldCounts(t *testing.T) {
	p, _ := setup(t, true)
	form, e := p.GetForm(context.Background())
	if e != nil {
		t.Fatal(e)
	}
	found := false
	for _, n := range form.Content {
		if n.Model == "cron" {
			found = n.Component == "VCronInput"
		}
	}
	if !found {
		t.Fatal("explicit cron component missing")
	}
	m, e := contracts.ParseManifest(ManifestJSON())
	if e != nil {
		t.Fatal(e)
	}
	scopes := []string{}
	for _, s := range m.Permissions {
		scopes = append(scopes, s.Scope)
	}
	for _, expr := range []string{"0 8 * * *", "0 0 8 * * *", "0 0 8 * * * 2030"} {
		p.config.Cron = expr
		services, e := p.GetService(context.Background())
		if e != nil {
			t.Fatal(e)
		}
		v := services[0]
		d := contracts.Descriptor{Kind: "service", ID: v.ID, Handler: contracts.HandlerID("service", v.ID), Name: v.Name, Trigger: v.Trigger, Cron: v.Cron, ScheduleSetting: v.ScheduleSetting, RequiredScopes: append([]string{"services.register"}, v.RequiredScopes...)}
		s := contracts.Snapshot{Version: contracts.DescriptorVersion, Descriptors: []contracts.Descriptor{d}}
		if e := contracts.ValidateSnapshot(m, scopes, s); e != nil {
			t.Fatal(expr, e)
		}
		for _, bad := range []string{"* * * *", "* * * * * * * *", strings.Repeat("x", 129)} {
			s.Descriptors[0].Cron = bad
			if contracts.ValidateSnapshot(m, scopes, s) == nil {
				t.Fatal(bad)
			}
		}
	}
}
