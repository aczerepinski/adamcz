package blog

import (
	"fmt"
	"strings"
	"time"

	"github.com/aczerepinski/adamcz/src/markdown"
)

// Post represents a Blog or Music post
type Post struct {
	Title          string
	Date           time.Time
	Description    string
	Body           string
	Tags           []string
	YouTubeLink    string
	SoundCloudLink string
	Composers      []string
	Performers     []string
	Slug           string
}

// NewPost initializes a post from raw key/value pairs
func NewPost(raw map[string]string) *Post {
	date, _ := time.Parse("1/2/06", raw["date"])
	return &Post{
		Title:          raw["title"],
		Date:           date,
		Description:    raw["description"],
		Body:           markdown.ToHTML(raw["body"]),
		Tags:           strings.Split(raw["tags"], ", "),
		YouTubeLink:    strings.Replace(raw["youtube"], "/watch?v=", "/embed/", 1),
		SoundCloudLink: raw["soundcloud"],
		Composers:      cleanEmpty(strings.Split(raw["composers"], ", ")),
		Performers:     cleanEmpty(strings.Split(raw["performers"], ", ")),
		Slug:           strings.ReplaceAll(strings.ToLower(raw["title"]), " ", "-"),
	}
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

// Summary prepares a summary that is contextually appropriate for
// music or tech blog posts
func (p *Post) Summary() string {
	if len(p.Composers) > 0 {
		return fmt.Sprintf("Recorded %s, composed by %s, performed by %s",
			p.Date.Format("1/2/06"), strings.Join(p.Composers, ","),
			strings.Join(p.Performers, ","))
	}
	return fmt.Sprintf("%s - %s", p.Date.Format("1/2/06"), p.Description)
}

// func (p *Post) MusicCredits() string {
// 	return fmt.Sprintf("composed by %s, performed by %s",
// 			strings.Join(p.Composers, ","),
// 			strings.Join(p.Performers, ","))
// }
