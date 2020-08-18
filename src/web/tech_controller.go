package web

import (
	"fmt"
	"net/http"
)

// techIndex serves a summary of all music posts
func (c *Controller) techIndex(w http.ResponseWriter, r *http.Request) {
	for _, p := range c.techPosts.GetAll(1, 10) {
		fmt.Fprintf(w, "%+v", p)
	}
}

// techShow serves a single music post
func (c *Controller) techShow(slug string, w http.ResponseWriter, r *http.Request) {
	post, err := c.techPosts.GetBySlug(slug)
	if err != nil {
		c.notFound(w, r)
		return
	}

	fmt.Fprintf(w, "%+v", post)
}
