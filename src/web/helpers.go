package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aczerepinski/adamcz/src/blog"
)

func getParam(param string, r *http.Request) []string {
	r.ParseForm()
	var vals []string
	if v := r.Form.Get(param); v != "" {
		// todo: parse multiple values
		vals = []string{r.Form.Get(param)}
	}
	return vals
}

// musicPath builds a /music URL with the given instruments and composers params,
// omitting any param whose slice is empty.
func musicPath(instruments, composers []string) string {
	var parts []string
	if len(instruments) > 0 {
		parts = append(parts, fmt.Sprintf("instruments=%s", url.QueryEscape(instruments[0])))
	}
	if len(composers) > 0 {
		parts = append(parts, fmt.Sprintf("composers=%s", url.QueryEscape(composers[0])))
	}
	if len(parts) == 0 {
		return "/music"
	}
	return "/music?" + strings.Join(parts, "&")
}

func musicFilters(activeInstruments []string, activeComposers []string) []filter {
	all := []string{
		"trumpet", "flugelhorn", "piano", "organ", "rhodes", "bass", "drums", "flute",
	}

	var filters []filter
	for _, instrument := range all {
		active := false
		for _, ai := range activeInstruments {
			if ai == instrument {
				active = true
			}
		}
		var path string
		if active {
			// deselect this instrument, preserve composers
			path = musicPath(nil, activeComposers)
		} else {
			path = musicPath([]string{instrument}, activeComposers)
		}
		filters = append(filters, filter{
			Name:   instrument,
			Active: active,
			Path:   path,
		})
	}
	return filters
}

type composerButton struct {
	label string
	value string
}

func composerFilters(activeComposers []string, activeInstruments []string) []filter {
	buttons := []composerButton{
		{label: "originals", value: "Adam Czerepinski"},
		{label: "covers", value: blog.NegationOperator + "Adam Czerepinski"},
	}

	var filters []filter
	for _, btn := range buttons {
		active := len(activeComposers) > 0 && activeComposers[0] == btn.value
		var path string
		if active {
			// deselect this composer, preserve instruments
			path = musicPath(activeInstruments, nil)
		} else {
			path = musicPath(activeInstruments, []string{btn.value})
		}
		filters = append(filters, filter{
			Name:   btn.label,
			Active: active,
			Path:   path,
		})
	}
	return filters
}
