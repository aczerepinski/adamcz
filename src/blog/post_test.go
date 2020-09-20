package blog

import "testing"

func TestNewPost(t *testing.T) {
	raw := map[string]string{
		"title": "Example",
		"date":  "8/15/20",
		"tags":  "Examples, Other Stuff",
		"body":  "This is the example body",
	}

	post := NewPost(raw)

	if post.Title != "Example" {
		t.Errorf("expected title to be Example, got %s", post.Title)
	}

	if post.Date.Month() != 8 || post.Date.Day() != 15 {
		t.Errorf("expected date to be 8/15, got %d/%d", post.Date.Month(), post.Date.Day())
	}

	if len(post.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(post.Tags))
	}

	if post.Body != "<p>This is the example body</p>" {
		t.Errorf("expected body to be '<p>This is the example body</p>', got %s", post.Body)
	}
}
