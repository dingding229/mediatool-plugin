package contracts

import "testing"

func TestAccessGuidanceIsStatic(t *testing.T) {
	for _, r := range []string{"disabled", "read_only", "permission_review_required"} {
		e := SafeError(AccessError(r))
		if e.Data.Reason != r || e.Message == "permission_denied" {
			t.Fatal(e)
		}
	}
	e := SafeError(&Error{Code: PermissionDenied, Message: "secret", Data: ErrorData{Reason: "secret"}})
	if e.Message != "permission_denied" || e.Data.Reason != "permission_denied" {
		t.Fatal(e)
	}
}
