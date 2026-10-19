package hello

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"testing"
)

func TestToolDescriptorsGoldenAndPureEcho(t *testing.T) {
	tools, e := (&Plugin{}).GetAgentTools(context.Background())
	if e != nil || len(tools) != 2 {
		t.Fatal(e)
	}
	for _, v := range tools {
		if _, e := contracts.ParseToolSchema(v.InputSchema); e != nil {
			t.Fatal(e)
		}
	}
	out, e := tools[0].Handler(context.Background(), json.RawMessage(`{"message":"golden"}`))
	if e != nil || string(out) != `{"message":"golden"}` {
		t.Fatal(string(out), e)
	}
}
