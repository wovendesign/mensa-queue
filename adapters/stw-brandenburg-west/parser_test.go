package stw_brandenburg_west

import "testing"

func TestSourceID(t *testing.T) {
	for _, tt := range []struct {
		dishID int64
		want   string
	}{
		{dishID: 0, want: "swt:0"},
		{dishID: 4711, want: "swt:4711"},
		{dishID: -1, want: "swt:-1"},
	} {
		if got := sourceID(tt.dishID); got != tt.want {
			t.Errorf("sourceID(%d) = %q, want %q", tt.dishID, got, tt.want)
		}
	}
}
