package tests

import (
	"reflect"
	"testing"
	"wardrobe/utils"
)

func TestCheckIfEmpty(t *testing.T) {
	type testCase[T any] struct {
		name     string
		input    []T
		expected any
	}

	// Sample Test
	tests := []testCase[int]{
		{
			name:     "Non-empty slice",
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Empty slice",
			input:    []int{},
			expected: nil,
		},
	}

	// Exec
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CheckIfEmpty(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CheckIfEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}
