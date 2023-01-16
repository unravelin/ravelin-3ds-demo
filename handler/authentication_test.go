package handler

import (
	"fmt"
	"testing"
)

func Test_convertToValidColorDepth(t *testing.T) {
	tests := []struct {
		colorDepth int
		expected   int
		error      bool
	}{
		{colorDepth: -10, error: true},
		{colorDepth: 0, error: true},
		{colorDepth: 1, expected: 1},
		{colorDepth: 3, expected: 1},
		{colorDepth: 4, expected: 4},
		{colorDepth: 6, expected: 4},
		{colorDepth: 8, expected: 8},
		{colorDepth: 13, expected: 8},
		{colorDepth: 15, expected: 15},
		{colorDepth: 16, expected: 16},
		{colorDepth: 20, expected: 16},
		{colorDepth: 24, expected: 24},
		{colorDepth: 30, expected: 24},
		{colorDepth: 32, expected: 32},
		{colorDepth: 40, expected: 32},
		{colorDepth: 48, expected: 48},
		{colorDepth: 50, expected: 48},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("convertToValidColorDepth(%d)", tt.colorDepth), func(t *testing.T) {
			actual, err := convertToValidColorDepth(tt.colorDepth)
			if err != nil {
				if tt.error {
					return
				}
				t.Fatalf("expected nil error, actual: %v", err)
			}

			if tt.error {
				t.Fatal("expected non nil error")
			}

			if actual != tt.expected {
				t.Fatalf("expected: %d, actual: %d", tt.expected, actual)
			}
		})
	}
}
