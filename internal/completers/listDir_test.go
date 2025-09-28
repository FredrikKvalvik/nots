package completers

import (
	"testing"

	"github.com/fredrikkvalvik/nots/internal/lister"
)

func AssertEqual(t *testing.T, got, want any, expect string) {
	if want != got {
		t.Fatalf("%s: want=%v got=%v", expect, want, got)
	}
}

func AssertNil(t *testing.T, val any, expect string) {
	if val != nil {
		t.Fatalf("%s: val=%v", expect, val)
	}
}
func AssertTrue(t *testing.T, val bool, expect string) {
	if !val {
		t.Fatalf("%s: val=%v", expect, val)
	}
}

func TestListDir(t *testing.T) {
	tests := []struct {
		ToComplete string
		Input      []lister.Path
		Out        []string
	}{
		{
			ToComplete: "",
			Input: []lister.Path{
				{"file.md"},
				{"dir", "file.md"},
				{"dir2"},
			},
			Out: []string{
				"dir",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.ToComplete, func(t *testing.T) {
			res := directoryList(tt.Input, tt.ToComplete)

			t.Log(res)
			AssertEqual(t, len(res), len(tt.Out), "results should equal expected length")

			for idx, line := range res {
				AssertEqual(t, line, tt.Out[idx], "each line should be equal")
			}
		})
	}
}
