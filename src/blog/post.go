package blog

import (
	"strings"
	"time"
)

// Post represents a Blog or Music post
type Post struct {
	Title          string
	Date           time.Time
	Body           string
	Tags           []string
	YouTubeLink    string
	SoundCloudLink string
	Composers      []string
	Performers     []string
	slug           string
}

// NewPost initializes a post from raw key/value pairs
func NewPost(raw map[string]string) *Post {
	date, _ := time.Parse("1/2/06", raw["date"])
	return &Post{
		Title:          raw["title"],
		Date:           date,
		Body:           raw["body"],
		Tags:           strings.Split(raw["tags"], ", "),
		YouTubeLink:    raw["youtube"],
		SoundCloudLink: raw["soundcloud"],
		Composers:      strings.Split(raw["composers"], ", "),
		Performers:     strings.Split(raw["performers"], ", "),
		slug:           strings.ReplaceAll(strings.ToLower(raw["title"]), " ", "-"),
	}
}
