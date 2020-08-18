package main

import (
	"log"
	"net/http"

	"github.com/aczerepinski/adamcz/src/blog"
	"github.com/aczerepinski/adamcz/src/web"
)

func main() {
	techPosts, err := blog.InitializeRepository("./techPosts")
	if err != nil {
		log.Fatalf("no tech posts! %v", err)
	}

	musicPosts, err := blog.InitializeRepository("./musicPosts")
	if err != nil {
		log.Fatalf("no music posts! %v", err)
	}

	controller := web.NewController(techPosts, musicPosts)
	http.Handle("/", controller)
	port := ":3000"
	http.ListenAndServe(port, nil)
}
