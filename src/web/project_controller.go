package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aczerepinski/adamcz/src/project"
)

type projectShow struct {
	CDN       string
	Version   string
	Project   project.Project
	PageTitle string
	MetaTitle string
}

func (c *Controller) projectRouter(w http.ResponseWriter, r *http.Request) {
	var slug string
	slug, r.URL.Path = popFromPath(r.URL.Path)
	c.projectShow(slug, w, r)
}

func (c *Controller) projectShow(slug string, w http.ResponseWriter, r *http.Request) {
	project := c.projects[slug]
	c.templates["projectShow"].Execute(w, projectShow{
		CDN:       cdnHost,
		Version:   c.version,
		Project:   project,
		PageTitle: project.Name,
		MetaTitle: fmt.Sprintf("adamcz | %s", strings.ToLower(project.Name)),
	})
}
