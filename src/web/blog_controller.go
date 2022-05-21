package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aczerepinski/adamcz/src/blog"
)

type blogIndex struct {
	CDN        string
	Version    string
	PageTitle  string
	MetaTitle  string
	PathPrefix string
	Posts      []*blog.Post
}

type blogShow struct {
	CDN          string
	Version      string
	PageTitle    string
	MetaTitle    string
	Post         *blog.Post
	PathPrefix   string
	MoreLikeThis []*blog.Post
}

// musicIndex serves a summary of all music posts
func (c *Controller) musicIndex(w http.ResponseWriter, r *http.Request) {
	data := blogIndex{
		CDN:        cdnHost,
		Version:    c.version,
		PageTitle:  "Music",
		MetaTitle:  "adamcz | music",
		PathPrefix: "/music/",
		Posts:      c.musicPosts.GetAll(1, 10),
	}
	c.templates["blogIndex"].Execute(w, data)
}

// techIndex serves a summary of all music posts
func (c *Controller) techIndex(w http.ResponseWriter, r *http.Request) {
	data := blogIndex{
		CDN:        cdnHost,
		Version:    c.version,
		PageTitle:  "Tech Blog",
		MetaTitle:  "adamcz | blog",
		PathPrefix: "/blog/",
		Posts:      c.techPosts.GetAll(1, 10),
	}
	c.templates["blogIndex"].Execute(w, data)
}

// musicShow serves a single music post
func (c *Controller) musicShow(slug string, w http.ResponseWriter, r *http.Request) {
	post, err := c.musicPosts.GetBySlug(slug)
	if err != nil {
		c.notFound(w, r)
		return
	}

	data := blogShow{
		Version:      c.version,
		PageTitle:    post.Title,
		MetaTitle:    fmt.Sprintf("adamcz | %s", strings.ToLower(post.Title)),
		Post:         post,
		PathPrefix:   "/music/",
		MoreLikeThis: c.musicPosts.GetRelateds(post, 3),
	}

	c.templates["blogShow"].Execute(w, data)
}

// techShow serves a single music post
func (c *Controller) techShow(slug string, w http.ResponseWriter, r *http.Request) {
	post, err := c.techPosts.GetBySlug(slug)
	if err != nil {
		c.notFound(w, r)
		return
	}

	data := blogShow{
		Version:      c.version,
		PageTitle:    post.Title,
		MetaTitle:    fmt.Sprintf("adamcz | %s", strings.ToLower(post.Title)),
		Post:         post,
		PathPrefix:   "/blog/",
		MoreLikeThis: c.techPosts.GetRelateds(post, 3),
	}

	c.templates["blogShow"].Execute(w, data)
}
