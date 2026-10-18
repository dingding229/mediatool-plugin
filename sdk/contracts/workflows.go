package contracts

import (
	"encoding/json"
	"regexp"
)

type WorkflowStep struct {
	ID         string        `json:"id"`
	InstanceID string        `json:"instanceId"`
	Action     string        `json:"action"`
	DependsOn  []string      `json:"dependsOn"`
	Input      ActionContext `json:"input"`
}
type Workflow struct {
	ID    string         `json:"id"`
	Steps []WorkflowStep `json:"steps"`
}
type WorkflowRequest struct {
	WorkflowID   string `json:"workflowId"`
	RunID        string `json:"runId,omitempty"`
	OperationKey string `json:"operationKey,omitempty"`
}
type WorkflowRun struct {
	ID         string `json:"id"`
	WorkflowID string `json:"workflowId"`
	State      string `json:"state"`
	Step       string `json:"step,omitempty"`
}
type WorkflowStarted struct {
	RunID      string `json:"runId"`
	WorkflowID string `json:"workflowId"`
	Revision   int64  `json:"revision"`
}
type WorkflowBeforeStep struct {
	WorkflowID string `json:"workflowId"`
	StepID     string `json:"stepId"`
	Allow      bool   `json:"allow"`
}
type RestrictResult struct {
	Allow bool `json:"allow"`
}

// WorkflowOrder returns a stable topological order. Depth<=16, steps<=32;
// inputs are explicit safe values, never arbitrary host resource capabilities.
func WorkflowOrder(w Workflow) ([]WorkflowStep, error) {
	valid := regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)
	raw, _ := json.Marshal(w)
	if !valid.MatchString(w.ID) || len(w.Steps) == 0 || len(w.Steps) > 32 || len(raw) > 16<<10 {
		return nil, NewError(InvalidParams)
	}
	nodes := map[string]WorkflowStep{}
	for _, n := range w.Steps {
		if !valid.MatchString(n.ID) || nodes[n.ID].ID != "" || len(n.InstanceID) != 36 || n.Action == "" || len(n.Action) > 200 || n.Input.Data == nil {
			return nil, NewError(InvalidParams)
		}
		nodes[n.ID] = n
	}
	depths := map[string]int{}
	visiting := map[string]bool{}
	out := []WorkflowStep{}
	var visit func(string) (int, error)
	visit = func(id string) (int, error) {
		if d := depths[id]; d > 0 {
			return d, nil
		}
		n, ok := nodes[id]
		if !ok || visiting[id] {
			return 0, NewError(InvalidParams)
		}
		visiting[id] = true
		depth := 1
		seen := map[string]bool{}
		for _, dep := range n.DependsOn {
			if seen[dep] {
				return 0, NewError(InvalidParams)
			}
			seen[dep] = true
			d, e := visit(dep)
			if e != nil {
				return 0, e
			}
			depth = max(depth, d+1)
		}
		if depth > 16 {
			return 0, NewError(InvalidParams)
		}
		visiting[id] = false
		depths[id] = depth
		out = append(out, n)
		return depth, nil
	}
	for _, n := range w.Steps {
		if _, e := visit(n.ID); e != nil {
			return nil, e
		}
	}
	return out, nil
}
