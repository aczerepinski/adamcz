package blog

import "testing"

func TestGetBy(t *testing.T) {
	tests := []struct {
		tags        []string
		expectedLen int
	}{
		{[]string{"trumpet"}, 3},
		{[]string{"piano"}, 2},
		{[]string{"trumpet", "piano"}, 1},
	}

	repo := testRepo()
	for i, test := range tests {
		posts := repo.GetBy(test.tags)
		if len(posts) != test.expectedLen {
			t.Errorf("test %d: expected %d posts, got %d", i, test.expectedLen, len(posts))
		}
	}
}

func TestGetRelateds(t *testing.T) {
	repo := testRepo()
	relateds := repo.GetRelateds(repo.posts[0], 1)
	if len(relateds) != 1 {
		t.Errorf("expected 1 related post, got %d", len(relateds))
	}
	if relateds[0].Title != "trumpet 2" {
		t.Errorf("expected related to be trumpet 2 because they have the same tags, got '%s'", relateds[0].Title)
	}
}

func testRepo() Repository {
	return Repository{
		posts: []*Post{
			{
				Title: "trumpet",
				Tags:  []string{"trumpet"},
			},
			{
				Title: "trumpet & piano",
				Tags:  []string{"trumpet", "piano"},
			},
			{
				Title: "piano",
				Tags:  []string{"piano"},
			},
			{
				Title: "trumpet 2",
				Tags:  []string{"trumpet"},
			},
		},
	}
}
