package main

import (
	"embed"
	"log"

	"github.com/warjiang/page-spy-api/config"
	"github.com/warjiang/page-spy-api/container"
	"github.com/warjiang/page-spy-api/serve"
)

//go:embed dist/*
var publicContent embed.FS

func main() {
	container := container.Container()
	err := container.Provide(func() *config.StaticConfig {
		return &config.StaticConfig{
			DirName: "dist",
			Files:   publicContent,
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	serve.Run()
}
