package business_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	classification "github.com/dingding229/mediatool-plugin/sdk/host/classification"
	files "github.com/dingding229/mediatool-plugin/sdk/host/files"
	library "github.com/dingding229/mediatool-plugin/sdk/host/library"
	metadata "github.com/dingding229/mediatool-plugin/sdk/host/metadata"
	storage "github.com/dingding229/mediatool-plugin/sdk/host/storage"
	strm "github.com/dingding229/mediatool-plugin/sdk/host/strm"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"strings"
	"testing"
)

func TestTypedMediaClients(t *testing.T) {
	ctx := context.Background()
	op := b.Operation{Key: "media-test-operation", Confirm: true}
	_ = op
	t.Run("library.assets", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := library.New(h).Assets(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.assets")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.assets", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).Assets(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.candidates", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.LibraryPage
		h := host.New(caller{m})
		_, e := library.New(h).Candidates(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.candidates")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.candidates", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).Candidates(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.history", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.LibraryPage
		h := host.New(caller{m})
		_, e := library.New(h).History(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.history")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.history", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).History(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.scans", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := library.New(h).Scans(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.scans")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.scans", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).Scans(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.recognize", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := library.New(h).Recognize(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.recognize")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.recognize", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).Recognize(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.suggest", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := library.New(h).Suggest(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.suggest")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.suggest", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).Suggest(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.organize.preview", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := library.New(h).OrganizePreview(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.organize.preview")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.organize.preview", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).OrganizePreview(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.correct.preview", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ConfirmIdentity
		h := host.New(caller{m})
		_, e := library.New(h).CorrectPreview(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.correct.preview")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.correct.preview", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).CorrectPreview(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.scan.start", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ScanStart
		h := host.New(caller{m})
		_, e := library.New(h).StartScan(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.scan.start")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.scan.start", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).StartScan(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.scan.control", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Control
		h := host.New(caller{m})
		_, e := library.New(h).ControlScan(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.scan.control")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.scan.control", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).ControlScan(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.organize.confirm", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ConfirmIdentity
		h := host.New(caller{m})
		_, e := library.New(h).Organize(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.organize.confirm")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.organize.confirm", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).Organize(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.correct.apply", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ConfirmIdentity
		h := host.New(caller{m})
		_, e := library.New(h).Correct(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.correct.apply")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.correct.apply", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).Correct(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.refresh", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := library.New(h).Refresh(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.refresh")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.refresh", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).Refresh(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("library.history.hide", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.HistorySelections
		h := host.New(caller{m})
		_, e := library.New(h).HideHistory(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("library.history.hide")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("library.history.hide", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = library.New(h).HideHistory(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("fs.media.roots", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := files.New(h).Roots(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("fs.media.roots")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("fs.media.roots", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = files.New(h).Roots(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("fs.media.list", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileQuery
		h := host.New(caller{m})
		_, e := files.New(h).List(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("fs.media.list")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("fs.media.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = files.New(h).List(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("fs.media.stat", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileQuery
		h := host.New(caller{m})
		_, e := files.New(h).Stat(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("fs.media.stat")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("fs.media.stat", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = files.New(h).Stat(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("fs.media.get", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := files.New(h).Get(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("fs.media.get")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("fs.media.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = files.New(h).Get(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("fs.media.preview", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileInput
		h := host.New(caller{m})
		_, e := files.New(h).Preview(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("fs.media.preview")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("fs.media.preview", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = files.New(h).Preview(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("fs.media.apply", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Approval
		h := host.New(caller{m})
		_, e := files.New(h).Apply(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("fs.media.apply")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("fs.media.apply", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = files.New(h).Apply(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("fs.media.previewBatch", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileInputs
		h := host.New(caller{m})
		_, e := files.New(h).PreviewBatch(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("fs.media.previewBatch")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("fs.media.previewBatch", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = files.New(h).PreviewBatch(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("fs.media.approveBatch", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Approval
		h := host.New(caller{m})
		_, e := files.New(h).ApproveBatch(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("fs.media.approveBatch")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("fs.media.approveBatch", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = files.New(h).ApproveBatch(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("fs.media.readText", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileQuery
		h := host.New(caller{m})
		_, e := files.New(h).ReadText(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("fs.media.readText")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("fs.media.readText", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = files.New(h).ReadText(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.roots", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := storage.New(h).Roots(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.roots")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.roots", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).Roots(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.list", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileQuery
		h := host.New(caller{m})
		_, e := storage.New(h).List(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.list")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).List(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.stat", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileQuery
		h := host.New(caller{m})
		_, e := storage.New(h).Stat(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.stat")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.stat", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).Stat(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.get", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := storage.New(h).Get(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.get")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).Get(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.preview", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileInput
		h := host.New(caller{m})
		_, e := storage.New(h).Preview(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.preview")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.preview", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).Preview(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.apply", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Approval
		h := host.New(caller{m})
		_, e := storage.New(h).Apply(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.apply")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.apply", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).Apply(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.previewBatch", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileInputs
		h := host.New(caller{m})
		_, e := storage.New(h).PreviewBatch(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.previewBatch")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.previewBatch", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).PreviewBatch(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.approveBatch", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Approval
		h := host.New(caller{m})
		_, e := storage.New(h).ApproveBatch(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.approveBatch")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.approveBatch", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).ApproveBatch(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.mkdir", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileInput
		h := host.New(caller{m})
		_, e := storage.New(h).Mkdir(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.mkdir")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.mkdir", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).Mkdir(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("storage.relocate", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.FileInput
		h := host.New(caller{m})
		_, e := storage.New(h).Relocate(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("storage.relocate")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("storage.relocate", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = storage.New(h).Relocate(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("strm.entries", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.PageRequest
		h := host.New(caller{m})
		_, e := strm.New(h).Entries(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("strm.entries")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("strm.entries", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = strm.New(h).Entries(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("strm.diagnostics", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := strm.New(h).Diagnostics(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("strm.diagnostics")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("strm.diagnostics", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = strm.New(h).Diagnostics(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("strm.repair.preview", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.IDs
		h := host.New(caller{m})
		_, e := strm.New(h).Preview(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("strm.repair.preview")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("strm.repair.preview", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = strm.New(h).Preview(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("strm.repair.approve", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Approval
		h := host.New(caller{m})
		_, e := strm.New(h).Approve(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("strm.repair.approve")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("strm.repair.approve", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = strm.New(h).Approve(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("strm.repair.retry", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Approval
		h := host.New(caller{m})
		_, e := strm.New(h).Retry(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("strm.repair.retry")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("strm.repair.retry", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = strm.New(h).Retry(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.parse", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Text
		h := host.New(caller{m})
		_, e := metadata.New(h).Parse(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.parse")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.parse", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).Parse(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.attributes.parse", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Text
		h := host.New(caller{m})
		_, e := metadata.New(h).Attributes(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.attributes.parse")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.attributes.parse", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).Attributes(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.attributes.validate", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Attribute
		h := host.New(caller{m})
		_, e := metadata.New(h).ValidateAttributes(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.attributes.validate")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.attributes.validate", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).ValidateAttributes(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.name.render", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.RenderName
		h := host.New(caller{m})
		_, e := metadata.New(h).Render(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.name.render")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.name.render", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).Render(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.episodes", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Text
		h := host.New(caller{m})
		_, e := metadata.New(h).Episodes(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.episodes")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.episodes", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).Episodes(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.normalize", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.EpisodeCoverage
		h := host.New(caller{m})
		_, e := metadata.New(h).Normalize(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.normalize")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.normalize", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).Normalize(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.search", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.MetadataQuery
		h := host.New(caller{m})
		_, e := metadata.New(h).Search(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.search")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.search", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).Search(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.details", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.MetadataQuery
		h := host.New(caller{m})
		_, e := metadata.New(h).Details(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.details")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.details", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).Details(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.storedDetails", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.MetadataQuery
		h := host.New(caller{m})
		_, e := metadata.New(h).StoredDetails(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.storedDetails")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.storedDetails", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).StoredDetails(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("metadata.refresh", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.ID
		h := host.New(caller{m})
		_, e := metadata.New(h).Refresh(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("metadata.refresh")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("metadata.refresh", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = metadata.New(h).Refresh(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("classification.get", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Empty
		h := host.New(caller{m})
		_, e := classification.New(h).Get(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("classification.get")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("classification.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = classification.New(h).Get(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("classification.match", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Classify
		h := host.New(caller{m})
		_, e := classification.New(h).Match(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("classification.match")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("classification.match", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = classification.New(h).Match(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("classification.render", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Classify
		h := host.New(caller{m})
		_, e := classification.New(h).Render(ctx, in)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("classification.render")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("classification.render", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = classification.New(h).Render(ctx, in)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
	t.Run("classification.update", func(t *testing.T) {
		m := plugintest.NewHost(nil)
		var in b.Classification
		h := host.New(caller{m})
		_, e := classification.New(h).Update(ctx, in, op)
		var rpc *contracts.Error
		if !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
			t.Fatal(e)
		}
		spec, _ := b.Lookup("classification.update")
		m = plugintest.NewHost([]string{spec.Scope})
		h = host.New(caller{m})
		_ = m.OnBusiness("classification.update", func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return nil, contracts.NewError(contracts.Unsupported)
		})
		_, e = classification.New(h).Update(ctx, in, op)
		if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
			t.Fatal(e)
		}
	})
}

func TestMediaDTOGolden(t *testing.T) {
	plugintest.Golden(t, "testdata/media-dto.json", struct {
		File     b.MediaFile    `json:"file"`
		Identity b.IdentityPlan `json:"identity"`
		Approval b.Approval     `json:"approval"`
	}{File: b.MediaFile{ID: "opaque-file", Name: "Owned.mkv", Path: "Movies/Owned.mkv", Size: 42}, Identity: b.IdentityPlan{ID: "opaque-asset", Revision: 3, PlanDigest: "digest", Operations: []b.IdentityOperation{{Path: "Movies/Owned.strm", Kind: "strm", Operation: "create"}}}, Approval: b.Approval{ID: "opaque-plan", Digest: "digest"}})
}

func TestMonitorObservationDTOIncludesOnlyBoundedIdentityEvidence(t *testing.T) {
	m := b.MediaFile{ID: "inode", Path: "incoming/Example.2026.mkv", Name: "Example.2026.mkv", Size: 4, Modified: 1234}
	raw, e := json.Marshal(m)
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(string(raw), `"modified":1234`) {
		t.Fatal(string(raw))
	}
	input := b.FileInput{RootID: "root", Operation: "organize", Source: m.Path, Destination: "organized", ExpectedFileID: m.ID, ExpectedSize: &m.Size, ExpectedModified: &m.Modified}
	raw, e = json.Marshal(input)
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(string(raw), "token") || strings.Contains(string(raw), "http") {
		t.Fatal(string(raw))
	}
}
