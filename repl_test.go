package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "      hello          world      ",
			expected: []string{"hello", "world"},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: "electrode DIGLETT Nidoran mAnKeY",
			expected: []string{"electrode", "diglett", "nidoran", "mankey"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Input: %v\nExpecting: %v\nActual: %v\nFail", c.input, c.expected, actual)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Input: %v\nExpecting: %v\nActual: %v\nFail", c.input, c.expected, actual)
				break
			}
		}
	}
}