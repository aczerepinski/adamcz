package blog

import "time"

// Post represents a Blog or Music post
type Post struct {
	Title          string
	Date           time.Time
	Body           string
	Tags           []string
	YouTubeLink    string
	SoundCloudLink string
}
