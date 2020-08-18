package web

import (
	"fmt"
	"net/http"
)

// musicIndex serves a summary of all music posts
func (c *Controller) musicIndex(w http.ResponseWriter, r *http.Request) {
	for _, p := range c.musicPosts.GetAll(1, 10) {
		fmt.Fprintf(w, "%+v", p)
	}
}

// musicShow serves a single music post
func (c *Controller) musicShow(slug string, w http.ResponseWriter, r *http.Request) {
	post, err := c.musicPosts.GetBySlug(slug)
	if err != nil {
		c.notFound(w, r)
		return
	}

	fmt.Fprintf(w, "%+v", post)
}
