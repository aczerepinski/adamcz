package web

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/aczerepinski/adamcz/src/blog"
)

// Controller handles routed requests
type Controller struct {
	techPosts  *blog.Repository
	musicPosts *blog.Repository
	templates  map[string]*template.Template
	version    string
}

// NewController returns an initialized controller
func NewController(version string, techPosts, musicPosts *blog.Repository) *Controller {
	return &Controller{
		version:    version,
		techPosts:  techPosts,
		musicPosts: musicPosts,
		templates:  initTemplates(),
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
	default:
		c.pageHandler(head, w, r)
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

func (c *Controller) pageHandler(slug string, w http.ResponseWriter, r *http.Request) {
	// PageData provides variables for interpolation into templates
	type PageData struct {
		Version   string
		MetaTitle string
		CDN       string
	}

	switch slug {
	case "resume":
		c.templates["resume"].Execute(w, PageData{CDN: cdnHost, Version: c.version, MetaTitle: "adamcz | resume"})
	case "bio":
		c.templates["bio"].Execute(w, PageData{CDN: cdnHost, Version: c.version, MetaTitle: "adamcz | bio"})
	case "photos":
		c.templates["photos"].Execute(w, PageData{CDN: cdnHost, Version: c.version, MetaTitle: "adamcz | photos"})
	case "":
		c.templates["home"].Execute(w, PageData{CDN: cdnHost, Version: c.version, MetaTitle: "adamcz"})
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

func initTemplates() map[string]*template.Template {
	root := "./src/templates/"
	bio := template.Must(template.ParseFiles(
		fmt.Sprintf("%slayout.html", root), fmt.Sprintf("%sbio.html", root)))

	blogIndex := template.Must(template.ParseFiles(
		fmt.Sprintf("%slayout.html", root), fmt.Sprintf("%sblogIndex.html", root)))

	blogShow := template.Must(template.ParseFiles(
		fmt.Sprintf("%slayout.html", root), fmt.Sprintf("%sblogShow.html", root)))

	resume := template.Must(template.ParseFiles(
		fmt.Sprintf("%slayout.html", root), fmt.Sprintf("%sresume.html", root)))

	home := template.Must(template.ParseFiles(
		fmt.Sprintf("%slayout.html", root), fmt.Sprintf("%shome.html", root)))

	photos := template.Must(template.ParseFiles(
		fmt.Sprintf("%slayout.html", root), fmt.Sprintf("%sphotos.html", root)))

	return map[string]*template.Template{
		"bio":       bio,
		"blogIndex": blogIndex,
		"blogShow":  blogShow,
		"home":      home,
		"photos":    photos,
		"resume":    resume,
	}
}
