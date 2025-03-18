package blog

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aczerepinski/adamcz/src/markdown"
)

// Post represents a Blog or Music post
type Post struct {
	Title        string
	Date         time.Time
	Description  string
	Body         string
	Tags         []string
	YouTubeLink  string
	SoundCloudID string
	Composers    []string
	Performers   []string
	Slug         string
	FilePath     string
	Project      string
}

// NewPost initializes a post from raw key/value pairs
func NewPost(raw map[string]string) *Post {
	date, _ := time.Parse("1/2/06", raw["date"])
	return &Post{
		Title:        raw["title"],
		Date:         date,
		Description:  raw["description"],
		Body:         markdown.ToHTML(raw["body"]),
		Tags:         strings.Split(raw["tags"], ", "),
		YouTubeLink:  strings.Replace(raw["youtube"], "/watch?v=", "/embed/", 1),
		SoundCloudID: raw["soundcloud"],
		Composers:    cleanEmpty(strings.Split(raw["composers"], ", ")),
		Performers:   cleanEmpty(strings.Split(raw["performers"], ", ")),
		Slug:         prepareSlug(raw["title"]),
		FilePath:     raw["filepath"],
		Project:      raw["project"],
	}
}

// Summary prepares a summary that is contextually appropriate for
// music or tech blog posts
func (p *Post) Summary() string {
	if p.performedByComposer() {
		return fmt.Sprintf("%scomposed and performed by %s",
			p.formattedDate(), p.Composers[0])
	}

	if len(p.Composers) > 0 {
		return fmt.Sprintf("%scomposed by %s, performed by %s",
			p.formattedDate(), strings.Join(p.Composers, ", "),
			strings.Join(p.Performers, ", "))
	}
	return fmt.Sprintf("%s %s", p.formattedDate(), p.Description)
}

func (p *Post) formattedDate() string {
	if p.Date.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s - ", p.Date.Format("1/2/06"))
}

func (p *Post) performedByComposer() bool {
	return len(p.Composers) == 1 &&
		len(p.Performers) == 1 &&
		p.Composers[0] == p.Performers[0]
}

func cleanEmpty(ss []string) []string {
	var cleaned []string
	for _, s := range ss {
		if clean := strings.TrimSpace(s); clean != "" {
			cleaned = append(cleaned, clean)
		}
	}
	return cleaned
}

func prepareSlug(title string) string {
	justWords := regexp.MustCompile(`[^\w\s]`).ReplaceAllString(title, "")
	return strings.ReplaceAll(strings.ToLower(justWords), " ", "-")
}
