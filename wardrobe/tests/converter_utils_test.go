package tests

import (
	"testing"
	"wardrobe/utils"
)

func TestBoolToYesNo(t *testing.T) {
	// Sample Test
	tests := []struct {
		name     string
		input    bool
		expected string
	}{
		{
			name:     "Input true returns Yes",
			input:    true,
			expected: "Yes",
		},
		{
			name:     "Input false returns No",
			input:    false,
			expected: "No",
		},
	}

	// Exec
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.BoolToYesNo(tt.input)
			if result != tt.expected {
				t.Errorf("BoolToYesNo(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
