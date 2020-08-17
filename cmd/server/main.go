package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	getPosts()
}

// Post is a post
type Post struct {
	Title string
	Body  string
	Date  time.Time
	tags  []string
}

func getPosts() []Post {
	var posts []Post
	root := "./blog/"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatalf("unable to read files: %v", err)
	}

	for _, f := range files {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", root, f.Name()))
		if err != nil {
			fmt.Printf("unable to read file: %v\n", err)
			continue
		}
		post, _ := parsePost(data)
		posts = append(posts, post)
	}

	return posts
}

func parsePost(file []byte) (Post, error) {
	fmt.Printf("file: %s\n", file)
	return Post{}, nil
}
