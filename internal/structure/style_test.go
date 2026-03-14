package structure_test

import (
	"testing"

	"scaff/internal/structure"
)

func TestMinimalPackageStyle_Resolve(t *testing.T) {
	style := structure.MinimalPackageStyle()

	tests := []struct {
		component string
		want      string
	}{
		{"root", "."},
		{"source", "."},
		{"test", "."},
		{"ci", ".github/workflows"},
		{"meta", "."},
		{"unknown", "."},
	}

	for _, tc := range tests {
		t.Run(tc.component, func(t *testing.T) {
			got := style.Resolve(tc.component)
			if got != tc.want {
				t.Errorf("Resolve(%q) = %q, want %q", tc.component, got, tc.want)
			}
		})
	}
}
