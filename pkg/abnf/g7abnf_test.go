package abnf

import (
	"testing"
)

func TestTag(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "HEAD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Tag(); got == nil {
				t.Errorf("Tag() should not return nil")
			}
		})
	}
}
