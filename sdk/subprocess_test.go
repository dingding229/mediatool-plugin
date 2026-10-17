package sdk_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/registry"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"github.com/dingding229/mediatool-plugin/sdk/transport"
	"os"
	"os/exec"
	"testing"
	"time"
)

// Test-only worker, never a distributable executable or the G2 application subcommand.
func TestWorkerHelperProcess(t *testing.T) {
	if os.Getenv("MEDIATOOL_SDK_TEST_WORKER") != "1" {
		return
	}
	d, e := registry.Get("hello-go")
	if e != nil {
		os.Exit(2)
	}
	p, e := transport.New(&transport.Stream{Reader: os.Stdin, Writer: os.Stdout}, transport.Options{Prefix: "p"})
	if e != nil {
		os.Exit(3)
	}
	r, e := sdk.NewRuntime(p, d.Manifest, d.Factory, sdk.Identity{InstanceID: "test-instance", DefinitionFingerprint: "test-definition", CatalogBuild: "test-build"})
	if e != nil {
		os.Exit(4)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = r.Serve(ctx) // Peer shutdown is expected after the parent finishes the protocol.
	os.Exit(0)
}
func TestRealGoChildStdioGolden(t *testing.T) {
	d, m, _, q := fixture(t)
	_ = d
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exe, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	cmd := exec.CommandContext(ctx, exe, "-test.run=^TestWorkerHelperProcess$")
	cmd.Env = []string{"MEDIATOOL_SDK_TEST_WORKER=1"}
	stdin, e := cmd.StdinPipe()
	if e != nil {
		t.Fatal(e)
	}
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		t.Fatal(e)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if e = cmd.Start(); e != nil {
		t.Fatal(e)
	}
	peer, e := transport.New(&transport.Stream{Reader: stdout, Writer: stdin}, transport.Options{Prefix: "h"})
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() {
		peer.Close()
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	if e = peer.Handle("host.call", m.Handle); e != nil {
		t.Fatal(e)
	}
	if e = peer.Start(ctx); e != nil {
		t.Fatal(e)
	}
	var initialized contracts.Initialized
	if e = peer.Call(ctx, "plugin.initialize", q, &initialized); e != nil {
		t.Fatal(e)
	}
	if e = peer.Call(ctx, "plugin.init", contracts.Config{Value: json.RawMessage(`{}`), Revision: 1}, nil); e != nil {
		t.Fatal(e)
	}
	var state bool
	var form sdk.Form
	var page sdk.Page
	for _, call := range []struct {
		method    string
		args, out any
	}{{"plugin.get_state", nil, &state}, {"plugin.get_form", nil, &form}, {"plugin.get_page", sdk.PageRequest{Key: "overview"}, &page}} {
		if e = peer.Call(ctx, call.method, call.args, call.out); e != nil {
			t.Fatal(e)
		}
	}
	plugintest.Golden(t, "testdata/hello-lifecycle.json", struct {
		Initialize contracts.Initialized `json:"initialize"`
		State      bool                  `json:"state"`
		Form       sdk.Form              `json:"form"`
		Page       sdk.Page              `json:"page"`
		Calls      []string              `json:"hostCalls"`
	}{initialized, state, form, page, m.Calls()})
	if e = peer.Call(ctx, "plugin.close", nil, nil); e != nil {
		t.Fatal(e)
	}
	peer.Close()
	if e = cmd.Wait(); e != nil {
		t.Fatalf("test worker failed (%T); stderr bytes=%d", e, stderr.Len())
	}
	if stderr.Len() != 0 {
		t.Fatal("unexpected stderr")
	}
	t.Log("real Go child stdio golden executed (not skipped); G2 app process isolation is not tested here")
}
