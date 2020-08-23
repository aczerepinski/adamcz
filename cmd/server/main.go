package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	version := time.Now().Format("20060102150405")

	controller := web.NewController(version, techPosts, musicPosts)
	http.Handle("/", controller)

	port := ":3000"
	fmt.Printf("version %s deployed\nlistening on port %s", version, port)
	http.ListenAndServe(port, nil)
}
