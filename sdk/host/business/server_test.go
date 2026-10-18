package business_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	catalog "github.com/dingding229/mediatool-plugin/sdk/host/catalog"
	mediaprobe "github.com/dingding229/mediatool-plugin/sdk/host/mediaprobe"
	mediaserver "github.com/dingding229/mediatool-plugin/sdk/host/mediaserver"
	mediausers "github.com/dingding229/mediatool-plugin/sdk/host/mediausers"
	metadata "github.com/dingding229/mediatool-plugin/sdk/host/metadata"
	playback "github.com/dingding229/mediatool-plugin/sdk/host/playback"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
)

func TestTypedServerClients(t *testing.T) {
	ctx := context.Background()
	op := b.Operation{Key: "server-test-operation", Confirm: true}
	t.Run("mediaserver.providers", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := mediaserver.New(h).Providers(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaserver.providers", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaserver.New(h).Providers(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaserver.libraries", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerQuery
		h := host.New(caller{m})
		_, e := mediaserver.New(h).Libraries(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaserver.libraries", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaserver.New(h).Libraries(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaserver.items", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerQuery
		h := host.New(caller{m})
		_, e := mediaserver.New(h).Items(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaserver.items", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaserver.New(h).Items(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaserver.settings.get", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := mediaserver.New(h).Settings(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaserver.settings.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaserver.New(h).Settings(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaserver.refresh", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerRefresh
		h := host.New(caller{m})
		_, e := mediaserver.New(h).Refresh(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.refresh"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaserver.refresh", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaserver.New(h).Refresh(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaserver.settings.update", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerConfiguration
		h := host.New(caller{m})
		_, e := mediaserver.New(h).UpdateSettings(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.settings.write"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaserver.settings.update", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaserver.New(h).UpdateSettings(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("catalog.list", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerQuery
		h := host.New(caller{m})
		_, e := catalog.New(h).List(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("catalog.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = catalog.New(h).List(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("catalog.item", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerQuery
		h := host.New(caller{m})
		_, e := catalog.New(h).Item(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("catalog.item", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = catalog.New(h).Item(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("catalog.home", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerQuery
		h := host.New(caller{m})
		_, e := catalog.New(h).Home(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("catalog.home", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = catalog.New(h).Home(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("catalog.discover", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerQuery
		h := host.New(caller{m})
		_, e := catalog.New(h).Discover(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("catalog.discover", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = catalog.New(h).Discover(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("catalog.nextUp", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerQuery
		h := host.New(caller{m})
		_, e := catalog.New(h).NextUp(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("catalog.nextUp", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = catalog.New(h).NextUp(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("catalog.userData", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ServerQuery
		h := host.New(caller{m})
		_, e := catalog.New(h).UserData(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("catalog.userData", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = catalog.New(h).UserData(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediausers.list", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := mediausers.New(h).List(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.users.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediausers.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediausers.New(h).List(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediausers.get", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := mediausers.New(h).Get(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.users.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediausers.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediausers.New(h).Get(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediausers.devices", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Selection
		h := host.New(caller{m})
		_, e := mediausers.New(h).Devices(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.users.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediausers.devices", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediausers.New(h).Devices(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediausers.history", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Selection
		h := host.New(caller{m})
		_, e := mediausers.New(h).History(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.users.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediausers.history", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediausers.New(h).History(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediausers.update", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.UserPatch
		h := host.New(caller{m})
		_, e := mediausers.New(h).Update(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.users.write"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediausers.update", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediausers.New(h).Update(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediausers.revokeDevice", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.DeviceRevocation
		h := host.New(caller{m})
		_, e := mediausers.New(h).RevokeDevice(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.users.write"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediausers.revokeDevice", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediausers.New(h).RevokeDevice(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("playback.sessions", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := playback.New(h).Sessions(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("playback.sessions", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = playback.New(h).Sessions(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("playback.summary", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := playback.New(h).Summary(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("playback.summary", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = playback.New(h).Summary(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("playeraccess.summary", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := playback.New(h).AccessSummary(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("playeraccess.summary", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = playback.New(h).AccessSummary(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaserver.sessions", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := mediaserver.New(h).Sessions(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaserver.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaserver.sessions", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaserver.New(h).Sessions(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaprobe.status", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := mediaprobe.New(h).Status(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaprobe.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaprobe.status", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaprobe.New(h).Status(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaprobe.list", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ProcessingQuery
		h := host.New(caller{m})
		_, e := mediaprobe.New(h).List(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaprobe.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaprobe.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaprobe.New(h).List(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaprobe.retry", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.IDs
		h := host.New(caller{m})
		_, e := mediaprobe.New(h).Retry(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaprobe.retry"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaprobe.retry", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaprobe.New(h).Retry(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaprobe.stills.retry", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.IDs
		h := host.New(caller{m})
		_, e := mediaprobe.New(h).RetryStills(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaprobe.retry"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaprobe.stills.retry", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaprobe.New(h).RetryStills(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("mediaprobe.stills.gc", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := mediaprobe.New(h).GC(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"mediaprobe.retry"})
		h = host.New(caller{m})
		_ = m.OnBusiness("mediaprobe.stills.gc", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = mediaprobe.New(h).GC(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.nativeDetails", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.NativeTarget
		h := host.New(caller{m})
		_, e := metadata.New(h).NativeDetails(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"metadata.read"})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.nativeDetails", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).NativeDetails(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.refreshDetails", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.NativeTarget
		h := host.New(caller{m})
		_, e := metadata.New(h).RefreshDetails(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		m = plugintest.NewHost([]string{"metadata.refresh"})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.refreshDetails", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).RefreshDetails(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
}
func TestServerDTOBoundaryGolden(t *testing.T) {
	v := struct {
		Item       b.ServerItem  `json:"item"`
		ActorQuery b.ServerQuery `json:"query"`
		Patch      b.UserPatch   `json:"patch"`
	}{Item: b.ServerItem{ID: "opaque-asset", Title: "Owned", Year: "2026", MediaType: "movie"}, ActorQuery: b.ServerQuery{PageRequest: b.PageRequest{Page: 1, Size: 200}, Provider: "builtin", LibraryID: "g:movie"}, Patch: b.UserPatch{ID: "opaque-user", Revision: 3, ClearExpiry: true}}
	plugintest.Golden(t, "testdata/server-dto.json", v)
	for _, name := range []string{"playback.resolve", "playback.streamURLs", "playback.activate", "mediaserver.playback", "mediaserver.image", "emby.gateway", "emby.proxy", "emby.managedLink", "emby.metadataHeaders", "mediausers.login", "mediausers.session", "mediausers.passwordHash", "mediausers.token", "playeraccess.events", "playeraccess.policy.update"} {
		if _, ok := b.Lookup(name); ok {
			t.Fatal("forbidden method exposed", name)
		}
	}
}
