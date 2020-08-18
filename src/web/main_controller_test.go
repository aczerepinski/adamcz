package web

import "testing"

func TestPopFromPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input        string
		expectedHead string
		expectedTail string
	}{
		{"", "", ""},
		{"/", "", ""},
		{"/music", "music", ""},
		{"/music/pizza", "music", "/pizza"},
	}

	for i, test := range tests {
		head, tail := popFromPath(test.input)
		if head != test.expectedHead {
			t.Errorf("test %d: expected head to be %s, got %s",
				i, test.expectedHead, head)
		}
		if tail != test.expectedTail {
			t.Errorf("test %d: expected tail to be %s, got %s",
				i, test.expectedTail, tail)
		}
	}
}
