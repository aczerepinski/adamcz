package project

import "github.com/aczerepinski/adamcz/src/blog"

type Project struct {
	Name         string
	Slug         string
	LogoURL      string
	PrimaryPhoto ProjectPhoto
	Description  string
	YouTubeURLs  []string
}

type ProjectPhoto struct {
	Url          string
	Photographer string
}

func (p ProjectPhoto) Empty() bool {
	return p.Url == "" && p.Photographer == ""
}

func InitProjects(allVideos []*blog.Post) map[string]Project {
	projects := map[string]Project{
		"hornado": {
			Name:    "Hornado",
			Slug:    "hornado",
			LogoURL: "/static/images/hornado_logo.png",
			Description: "Hornado is Madison's most devastating jazz/groove band. " +
				"Featuring four of Wisconsin's most extreme horn players, Hornado lays down twisting grooves and whirling solos that will leave you spinning. " +
				"Hornado specializes in typhoon-jazz, windstorm-funk, and 8-bit NES classics.",
			YouTubeURLs: []string{},
		},
		"piano-trio": {
			Name: "Adam Czerepinski Piano Trio",
			Slug: "piano-trio",
			Description: "Multi-instrumentalist Adam Czerepinski's trio is his primary outlet for exploring his original instrument - the piano. " +
				"In addition to Adam's original compositions, the trio explores popular music from the past twenty years. " +
				"The band's songbook includes covers of OutKast, Thundercat, The Smile, Chris Stapleton, Hiatus Kaiyote, Sufjan Stevens, and more.",
			YouTubeURLs: []string{},
			PrimaryPhoto: ProjectPhoto{
				Url:          "/static/images/piano_trio_primary.jpg",
				Photographer: "Kristin Shafel",
			},
		},
	}

	for _, video := range allVideos {
		if project, ok := projects[video.Project]; ok {
			project.YouTubeURLs = append(project.YouTubeURLs, video.YouTubeLink)
			projects[video.Project] = project
		}
	}

	return projects
}
