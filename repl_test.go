package main

import (
	"testing"
	"reflect"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  hello   world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Hello World",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		result := cleanInput(c.input)
		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("cleanInput(%q) = %v, want %v", c.input, result, c.expected)
		}
	}
}
