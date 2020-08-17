package blog

import "testing"

func TestParse(t *testing.T) {
	data := `*** title ***
	Example

	*** date ***
	8/15/20

	*** tags ***
	Examples, Other Stuff

	*** body ***
	This is the example body`

	parser := Parser{}

	parsed := parser.Parse([]byte(data))

	if title := parsed["title"]; title != "Example" {
		t.Errorf("expected title to be Example, got %s", title)
	}

	if body := parsed["body"]; body != "This is the example body" {
		t.Errorf("expected body to be This is the example body, got %s", body)
	}
}
