package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aczerepinski/adamcz/src/blog"
	"github.com/aczerepinski/adamcz/src/calendar"
	"github.com/aczerepinski/adamcz/src/project"
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

	transcriptions, err := blog.InitializeRepository("./transcriptions")
	if err != nil {
		log.Fatalf("no transcriptions! %v", err)
	}

	events, err := calendar.InitializeRepository("./events")
	if err != nil {
		log.Fatalf("no events! %v", err)
	}

	projects := project.InitProjects(musicPosts.GetAll(1, 10))

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	version := time.Now().Format("20060102150405")

	controller := web.NewController(version, techPosts, musicPosts, transcriptions, events, projects)
	http.Handle("/", controller)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("version %s deployed\nlistening on port %s", version, port)
	http.ListenAndServe(":"+port, nil)
}
