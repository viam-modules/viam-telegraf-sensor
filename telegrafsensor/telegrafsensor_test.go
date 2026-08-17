package telegrafsensor

import "testing"

func TestModelName(t *testing.T) {
	if Model.Name != "telegrafsensor" {
		t.Errorf("expected model name %q, got %q", "telegrafsensor", Model.Name)
	}
}
