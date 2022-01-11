package serialdet

import "testing"

func TestError(t *testing.T) {
	const msg = "test error"
	const err Error = Error(msg)

	if err.Error() != msg {
		t.Fatal("given", err.Error(), "Want", msg)
	}
}
