package blog

import (
	"fmt"
	"io/ioutil"
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

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return &r, fmt.Errorf("unable to read directory: %v", err)
	}

	for _, f := range files {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", root, f.Name()))
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
	var results []*Post
	for _, requestedTag := range p.Tags {
		for _, post := range r.posts {
			if p == post {
				continue
			}
			for _, tag := range post.Tags {
				if tag == requestedTag && !contains(results, post) {
					results = append(results, post)
					if len(results) == quantity {
						return results
					}
				}
			}
		}
	}

	return results
}

func contains(posts []*Post, p *Post) bool {
	for _, post := range posts {
		if p.Title == post.Title {
			return true
		}
	}
	return false
}
