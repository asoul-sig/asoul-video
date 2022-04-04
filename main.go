package main

import (
	gocontext "context"
	"io/fs"
	"net/http"
	"os"

	"github.com/flamego/flamego"
	"github.com/robfig/cron/v3"
	log "unknwon.dev/clog/v2"

	"github.com/asoul-sig/asoul-video/frontend"
	"github.com/asoul-sig/asoul-video/internal/context"
	"github.com/asoul-sig/asoul-video/internal/db"
	"github.com/asoul-sig/asoul-video/internal/route"
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

	// Register cron task.
	ctx := gocontext.Background()
	if err := db.Videos.Refresh(ctx); err != nil {
		log.Error("Failed to refresh materialized views: %v", err)
	}
	c := cron.New()
	if _, err := c.AddFunc("@every 1h", func() {
		if err := db.Videos.Refresh(ctx); err != nil {
			log.Error("Failed to refresh materialized views: %v", err)
		}
		log.Trace("Refresh materialized views.")
	}); err != nil {
		log.Fatal("Failed to add cron function: %v", err)
	}
	c.Start()

	f := flamego.Classic()
	f.Use(func(ctx flamego.Context) {
		ctx.ResponseWriter().Header().Set("Access-Control-Allow-Methods", http.MethodGet)
		ctx.ResponseWriter().Header().Set("Access-Control-Max-Age", "600")
		ctx.ResponseWriter().Header().Set("Access-Control-Allow-Origin", "*")
	})
	f.Use(context.Contexter())

	fe, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatal("Failed to sub filesystem: %v", err)
	}
	f.Use(flamego.Static(flamego.StaticOptions{
		FileSystem: http.FS(fe),
	}))

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
		f.Get("/video_urls", source.VideoURLs)
		f.Get("/video_ids", source.VideoIDs)
	}, source.VerifyKey(os.Getenv("SOURCE_REPORT_KEY")))

	f.Get("/ping", func(ctx context.Context) {
		ctx.Success(map[string]interface{}{
			"build_time":   BuildTime,
			"build_commit": BuildCommit,
		})
	})

	f.Run()
}
