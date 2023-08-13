package blog

import "testing"

func TestGetBy(t *testing.T) {
	tests := []struct {
		tags        []string
		expectedLen int
	}{
		{[]string{"trumpet"}, 2},
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

func testRepo() Repository {
	return Repository{
		posts: []*Post{
			&Post{
				Title: "trumpet",
				Tags:  []string{"trumpet"},
			},
			&Post{
				Title: "trumpet & piano",
				Tags:  []string{"trumpet", "piano"},
			},
			&Post{
				Title: "piano",
				Tags:  []string{"piano"},
			},
		},
	}
}
