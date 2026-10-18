package contracts_test

import (
	"encoding/json"
	"fmt"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"strings"
	"testing"
)

func TestWorkflowDAGAndAuthGolden(t *testing.T) {
	n := contracts.WorkflowStep{ID: "one", InstanceID: aid, Action: "hello", DependsOn: []string{}, Input: contracts.ActionContext{Data: map[string]json.RawMessage{}}}
	w := contracts.Workflow{ID: "demo", Steps: []contracts.WorkflowStep{n}}
	if _, e := contracts.WorkflowOrder(w); e != nil {
		t.Fatal(e)
	}
	for _, bad := range []contracts.Workflow{{ID: "bad", Steps: []contracts.WorkflowStep{func() contracts.WorkflowStep { x := n; x.DependsOn = []string{"one"}; return x }()}}, {ID: "bad", Steps: make([]contracts.WorkflowStep, 33)}, {ID: "bad", Steps: []contracts.WorkflowStep{n, n}}} {
		if _, e := contracts.WorkflowOrder(bad); e == nil {
			t.Fatal("invalid DAG")
		}
	}
	plugintest.Golden(t, "testdata/workflows.json", map[string]any{"workflow": w, "beforeStep": contracts.WorkflowBeforeStep{WorkflowID: "demo", StepID: "one", Allow: true}, "started": contracts.WorkflowStarted{WorkflowID: "demo", RunID: "fact-v1", Revision: 1}, "auth": contracts.AuthRequest{ProviderID: "fake", Challenge: "bound"}})
}

func TestWorkflowDepthUnknownDependencyAndInputBounds(t *testing.T) {
	w := contracts.Workflow{ID: "depth"}
	for i := 0; i < 17; i++ {
		n := contracts.WorkflowStep{ID: fmt.Sprint(i), InstanceID: aid, Action: "hello", Input: contracts.ActionContext{Data: map[string]json.RawMessage{}}}
		if i > 0 {
			n.DependsOn = []string{fmt.Sprint(i - 1)}
		}
		w.Steps = append(w.Steps, n)
	}
	if _, e := contracts.WorkflowOrder(w); e == nil {
		t.Fatal("depth")
	}
	w.Steps = w.Steps[:1]
	w.Steps[0].DependsOn = []string{"missing"}
	if _, e := contracts.WorkflowOrder(w); e == nil {
		t.Fatal("missing")
	}
	w.Steps[0].DependsOn = nil
	w.Steps[0].Input.Data["large"] = json.RawMessage(`"` + strings.Repeat("x", 17000) + `"`)
	if _, e := contracts.WorkflowOrder(w); e == nil {
		t.Fatal("input limit")
	}
}
