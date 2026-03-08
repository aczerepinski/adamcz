package blog

import "testing"

func TestGetBy(t *testing.T) {
	tests := []struct {
		name        string
		query       Query
		expectedLen int
	}{
		{
			name:        "instrument only",
			query:       Query{Instruments: []string{"trumpet"}},
			expectedLen: 3,
		},
		{
			name:        "multiple instruments",
			query:       Query{Instruments: []string{"trumpet", "piano"}},
			expectedLen: 1,
		},
		{
			name:        "instrument excluded",
			query:       Query{Instruments: []string{"!trumpet"}},
			expectedLen: 3,
		},
		{
			name:        "composer required",
			query:       Query{Composers: []string{"Herbie Hancock"}},
			expectedLen: 2,
		},
		{
			name:        "composer excluded (covers)",
			query:       Query{Composers: []string{"!Adam Czerepinski"}},
			expectedLen: 3,
		},
		{
			name:        "instrument and composer required",
			query:       Query{Instruments: []string{"piano"}, Composers: []string{"Herbie Hancock"}},
			expectedLen: 1,
		},
		{
			name:        "instrument required and composer excluded",
			query:       Query{Instruments: []string{"trumpet"}, Composers: []string{"!Adam Czerepinski"}},
			expectedLen: 1,
		},
		{
			name:        "instrument excluded and composer required",
			query:       Query{Instruments: []string{"!trumpet"}, Composers: []string{"Herbie Hancock"}},
			expectedLen: 1,
		},
	}

	repo := testRepo()
	for _, test := range tests {
		posts := repo.GetBy(test.query)
		if len(posts) != test.expectedLen {
			t.Errorf("%s: expected %d posts, got %d", test.name, test.expectedLen, len(posts))
		}
	}
}

func TestGetRelateds(t *testing.T) {
	repo := testRepo()
	relateds := repo.GetRelateds(repo.posts[0], 1)
	if len(relateds) != 1 {
		t.Errorf("expected 1 related post, got %d", len(relateds))
	}
	if relateds[0].Title != "piano" {
		t.Errorf("expected related to be piano because they share the same composer (Adam Czerepinski), got '%s'", relateds[0].Title)
	}
}

func testRepo() Repository {
	return Repository{
		posts: []*Post{
			{
				Title:     "trumpet",
				Tags:      []string{"trumpet"},
				Composers: []string{"Adam Czerepinski"},
			},
			{
				Title:     "trumpet & piano",
				Tags:      []string{"trumpet", "piano"},
				Composers: []string{"Adam Czerepinski"},
			},
			{
				Title:     "piano",
				Tags:      []string{"piano"},
				Composers: []string{"Adam Czerepinski"},
			},
			{
				Title:     "trumpet 2",
				Tags:      []string{"trumpet"},
				Composers: []string{"Herbie Hancock"},
			},
			{
				Title:     "herbie piano",
				Tags:      []string{"piano"},
				Composers: []string{"Herbie Hancock"},
			},
			{
				Title:     "bass",
				Tags:      []string{"bass"},
				Composers: []string{"Wayne Shorter"},
			},
		},
	}
}
