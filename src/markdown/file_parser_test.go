package markdown

import "testing"

func TestParseFile(t *testing.T) {
	data := `*** title ***
	Example

	*** date ***
	8/15/20

	*** tags ***
	Examples, Other Stuff

	*** body ***
	This is the example body`

	parsed := ParseFile([]byte(data))

	if title := parsed["title"]; title != "Example" {
		t.Errorf("expected title to be Example, got %s", title)
	}

	if date := parsed["date"]; date != "8/15/20" {
		t.Errorf("expected date to be 8/15/20, got %s", date)
	}

	if tags := parsed["tags"]; tags != "Examples, Other Stuff" {
		t.Errorf("expected tags to be Examples, Other Stuff, got %s", tags)
	}

	if body := parsed["body"]; body != "This is the example body" {
		t.Errorf("expected body to be This is the example body, got %s", body)
	}
}
