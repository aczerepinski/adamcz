package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aczerepinski/adamcz/src/blog"
)

// Controller handles routed requests
type Controller struct {
	techPosts  *blog.Repository
	musicPosts *blog.Repository
}

// NewController returns an initialized controller
func NewController(techPosts, musicPosts *blog.Repository) *Controller {
	return &Controller{
		techPosts:  techPosts,
		musicPosts: musicPosts,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = popFromPath(r.URL.Path)
	switch head {
	case "music":
		c.musicRouter(w, r)
	case "blog":
		c.techRouter(w, r)
	}
}

func (c *Controller) musicRouter(w http.ResponseWriter, r *http.Request) {
	var slug string
	slug, r.URL.Path = popFromPath(r.URL.Path)
	switch slug {
	case "":
		c.musicIndex(w, r)
	default:
		c.musicShow(slug, w, r)
	}
}

func (c *Controller) techRouter(w http.ResponseWriter, r *http.Request) {
	var slug string
	slug, r.URL.Path = popFromPath(r.URL.Path)
	switch slug {
	case "":
		c.techIndex(w, r)
	default:
		c.techShow(slug, w, r)
	}
}

func (c *Controller) notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "404")
}

func popFromPath(p string) (head, tail string) {
	if strings.Index(p, "/") == 0 {
		p = strings.Trim(p, "/")
	}

	i := strings.Index(p, "/")

	if i == -1 {
		return p, ""
	}
	return p[:i], p[i:]
}
