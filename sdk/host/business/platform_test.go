package business_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	cache "github.com/dingding229/mediatool-plugin/sdk/host/cache"
	integrations "github.com/dingding229/mediatool-plugin/sdk/host/integrations"
	jobs "github.com/dingding229/mediatool-plugin/sdk/host/jobs"
	maintenance "github.com/dingding229/mediatool-plugin/sdk/host/maintenance"
	network "github.com/dingding229/mediatool-plugin/sdk/host/network"
	notifications "github.com/dingding229/mediatool-plugin/sdk/host/notifications"
	settings "github.com/dingding229/mediatool-plugin/sdk/host/settings"
	system "github.com/dingding229/mediatool-plugin/sdk/host/system"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
)

func TestTypedPlatformClients(t *testing.T) {
	ctx := context.Background()
	op := b.Operation{Key: "platform-operation", Confirm: true}
	t.Run("integrations.list", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := integrations.New(h).List(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"integrations.read"}, b.RequiredScopes("integrations.list")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("integrations.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = integrations.New(h).List(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("integrations.status", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := integrations.New(h).Status(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"integrations.read"}, b.RequiredScopes("integrations.status")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("integrations.status", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = integrations.New(h).Status(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("integrations.test", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := integrations.New(h).Test(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"integrations.read"}, b.RequiredScopes("integrations.test")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("integrations.test", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = integrations.New(h).Test(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("settings.get", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.SettingKeys
		h := host.New(caller{m})
		_, e := settings.New(h).Get(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"settings.read"}, b.RequiredScopes("settings.get")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("settings.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = settings.New(h).Get(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("settings.patch", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.SettingPatch
		h := host.New(caller{m})
		_, e := settings.New(h).Patch(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"settings.write"}, b.RequiredScopes("settings.patch")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("settings.patch", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = settings.New(h).Patch(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("cache.status", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.NamespaceQuery
		h := host.New(caller{m})
		_, e := cache.New(h).Status(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"cache.read"}, b.RequiredScopes("cache.status")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("cache.status", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = cache.New(h).Status(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("cache.list", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.NamespaceQuery
		h := host.New(caller{m})
		_, e := cache.New(h).List(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"cache.read"}, b.RequiredScopes("cache.list")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("cache.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = cache.New(h).List(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("cache.invalidate", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.CacheInvalidation
		h := host.New(caller{m})
		_, e := cache.New(h).Invalidate(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"cache.invalidate"}, b.RequiredScopes("cache.invalidate")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("cache.invalidate", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = cache.New(h).Invalidate(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("jobs.list", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.JobQuery
		h := host.New(caller{m})
		_, e := jobs.New(h).List(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"jobs.read"}, b.RequiredScopes("jobs.list")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("jobs.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = jobs.New(h).List(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("jobs.get", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := jobs.New(h).Get(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"jobs.read"}, b.RequiredScopes("jobs.get")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("jobs.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = jobs.New(h).Get(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("jobs.cancel", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := jobs.New(h).Cancel(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"jobs.cancel"}, b.RequiredScopes("jobs.cancel")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("jobs.cancel", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = jobs.New(h).Cancel(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("system.info", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := system.New(h).Info(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"system.read"}, b.RequiredScopes("system.info")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("system.info", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = system.New(h).Info(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("system.health", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := system.New(h).Health(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"system.read"}, b.RequiredScopes("system.health")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("system.health", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = system.New(h).Health(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("system.modules", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := system.New(h).Modules(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"system.read"}, b.RequiredScopes("system.modules")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("system.modules", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = system.New(h).Modules(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("maintenance.health", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := maintenance.New(h).Health(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"maintenance.read"}, b.RequiredScopes("maintenance.health")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("maintenance.health", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = maintenance.New(h).Health(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("maintenance.alerts", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := maintenance.New(h).Alerts(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"maintenance.read"}, b.RequiredScopes("maintenance.alerts")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("maintenance.alerts", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = maintenance.New(h).Alerts(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("notifications.publish", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Notice
		h := host.New(caller{m})
		_, e := notifications.New(h).Send(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"notifications.send"}, b.RequiredScopes("notifications.publish")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("notifications.publish", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = notifications.New(h).Send(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("notifications.deliveries", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := notifications.New(h).Deliveries(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"notifications.read"}, b.RequiredScopes("notifications.deliveries")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("notifications.deliveries", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = notifications.New(h).Deliveries(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("notifications.retry", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.DeliveryRetry
		h := host.New(caller{m})
		_, e := notifications.New(h).Retry(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"notifications.retry"}, b.RequiredScopes("notifications.retry")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("notifications.retry", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = notifications.New(h).Retry(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("networkproxy.status", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := network.New(h).ProxyStatus(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"integrations.read"}, b.RequiredScopes("networkproxy.status")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("networkproxy.status", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = network.New(h).ProxyStatus(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("network.http", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.HostRequest
		h := host.New(caller{m})
		_, e := network.New(h).Request(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost(append([]string{"network.proxy.use"}, b.RequiredScopes("network.http")...))
		h = host.New(caller{m})
		_ = m.OnBusiness("network.http", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = network.New(h).Request(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
}
