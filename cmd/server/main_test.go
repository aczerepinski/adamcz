package main

import (
	"testing"

	"github.com/aczerepinski/adamcz/src/blog"
	"github.com/aczerepinski/adamcz/src/calendar"
)

func TestMusicPosts(t *testing.T) {
	_, err := blog.InitializeRepository("../../musicPosts")
	if err != nil {
		t.Fatalf("error initializing music repo: %v", err)
	}
}

func TestCalendar(t *testing.T) {
	_, err := calendar.InitializeRepository("../../events")
	if err != nil {
		t.Fatalf("error initializing event repo: %v", err)
	}
}
