package main

import (
	"os"

	"github.com/flamego/flamego"
	"github.com/flamego/template"
	log "unknwon.dev/clog/v2"

	"github.com/asoul-video/asoul-video/assets/templates"
	"github.com/asoul-video/asoul-video/internal/context"
	"github.com/asoul-video/asoul-video/internal/db"
	"github.com/asoul-video/asoul-video/internal/route"
)

var (
	BuildTime   string
	BuildCommit string
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

	fs, err := template.EmbedFS(templates.FS, ".", []string{".html"})
	if err != nil {
		log.Fatal("Failed to embed template file system: %v", err)
	}
	f.Use(template.Templater(template.Options{FileSystem: fs}))

	landing := route.NewLandingHandler()
	f.Get("/", landing.Home)
	f.NotFound(landing.Home)

	f.Group("/api", func() {
		member := route.NewMemberHandler()
		f.Get("/members", member.List)
		f.Get("/member/{secUID}", member.GetBySecUID)

		video := route.NewVideoHandler()
		f.Get("/videos", video.List)
		f.Group("/video", func() {
			f.Get("/{id}", video.GetByID)
			f.Get("/random", video.Random)
		})
	})

	// Crawler report service.
	source := route.NewSourceHandler()
	f.Group("/source", func() {
		f.Post("/report", source.Report)
	}, source.VerifyKey(os.Getenv("SOURCE_REPORT_KEY")))

	f.Get("/ping", func(ctx context.Context) {
		ctx.Success(map[string]interface{}{
			"build_time":   BuildTime,
			"build_commit": BuildCommit,
		})
	})

	f.Run()
}
