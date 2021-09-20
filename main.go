package main

import (
	"os"

	"github.com/flamego/flamego"
	log "unknwon.dev/clog/v2"

	"github.com/asoul-video/asoul-video/internal/context"
	"github.com/asoul-video/asoul-video/internal/db"
	"github.com/asoul-video/asoul-video/internal/route"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}

	f := flamego.Classic()
	f.Use(context.Contexter())

	// Crawler report service.
	source := route.NewSourceHandler()
	f.Group("/source", func() {
		f.Post("/report", source.Report)
	}, source.VerifyKey(os.Getenv("SOURCE_REPORT_KEY")))

	f.Run()
}
