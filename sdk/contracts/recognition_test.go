package contracts

import (
	"strings"
	"testing"
)

func TestNameHintSchema(t *testing.T) {
	v := NameRecognition{Name: "Demo.mkv", Hint: NameHint{Title: "Demo", MediaType: "movie"}, Continue: true}
	if ValidateName(v) != nil {
		t.Fatal("valid")
	}
	for _, mutate := range []func(*NameRecognition){func(x *NameRecognition) { x.Name = "/secret/file" }, func(x *NameRecognition) { x.Hint.Title = strings.Repeat("x", 513) }, func(x *NameRecognition) { x.Hint.MediaType = "secret" }, func(x *NameRecognition) { x.Hint.Year = "abcd" }, func(x *NameRecognition) { x.Hint.Episodes = []int{-1} }} {
		x := v
		mutate(&x)
		if ValidateName(x) == nil {
			t.Fatal("unsafe", x)
		}
	}
}
