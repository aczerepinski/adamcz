package blog

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aczerepinski/adamcz/src/markdown"
)

// Repository gives read access to blog posts
type Repository struct {
	posts []*Post
}

// InitializeRepository initializes from a directory markdown
// files. For an example of the expected file format see the tests
// for ParseFile.
func InitializeRepository(root string) (*Repository, error) {
	var r Repository

	if !strings.HasSuffix("/", root) {
		root = fmt.Sprintf("%s/", root)
	}

	files, err := os.ReadDir(root)
	if err != nil {
		return &r, fmt.Errorf("unable to read directory: %v", err)
	}

	for _, f := range files {
		data, err := os.ReadFile(fmt.Sprintf("%s%s", root, f.Name()))
		if err != nil {
			fmt.Printf("unable to read file: %v\n", err)
			continue
		}

		raw := markdown.ParseFile(data)
		r.posts = append(r.posts, NewPost(raw))
	}

	if len(r.posts) == 0 {
		return &r, fmt.Errorf("sorry, no posts")
	}

	sort.Slice(r.posts, func(i, j int) bool {
		return r.posts[j].Date.Before(r.posts[i].Date)
	})

	return &r, nil
}

// GetAll returns all posts
func (r *Repository) GetAll(page, perPage int) []*Post {
	return r.posts
}

// GetBy returns posts that contain a matching tag for ALL query terms
func (r *Repository) GetBy(query []string) []*Post {
	var posts []*Post
outer:
	for _, post := range r.posts {
		for _, searchTerm := range query {
			found := false
			for _, postTag := range post.Tags {
				if postTag == searchTerm {
					found = true
				}
			}
			if !found {
				continue outer
			}
		}
		posts = append(posts, post)
	}
	return posts
}

// GetBySlug returns a single Post
func (r *Repository) GetBySlug(slug string) (*Post, error) {
	for _, p := range r.posts {
		if p.Slug == slug {
			return p, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

// GetRelateds returns posts containing overlapping tags
func (r *Repository) GetRelateds(p *Post, quantity int) []*Post {
	type match struct {
		score int
		post  *Post
	}
	var matches []match

	for _, post := range r.posts {
		if p == post {
			continue
		}
		score := similarityScore(p, post)
		if len(matches) < quantity {
			matches = append(matches, match{score, post})
			continue
		}
		for i, m := range matches {
			if score > m.score {
				matches[i] = match{score, post}
				break
			}
		}
	}

	var posts []*Post
	for _, m := range matches {
		posts = append(posts, m.post)
	}
	return posts
}

func contains(posts []*Post, p *Post) bool {
	for _, post := range posts {
		if p.Title == post.Title {
			return true
		}
	}
	return false
}

func similarityScore(a *Post, b *Post) int {
	var score int
	if len(a.Tags) == len(b.Tags) {
		score += 3
	}
	for _, tag := range a.Tags {
		for _, match := range b.Tags {
			if tag == match {
				score += 1
				break
			}
		}
	}
	if len(a.Composers) == 0 || len(b.Composers) == 0 {
		return score
	}
	if a.Composers[0] == b.Composers[0] {
		score += 7
	}
	return score
}
