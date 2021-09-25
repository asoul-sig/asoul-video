// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	log "unknwon.dev/clog/v2"

	"github.com/asoul-video/asoul-video/internal/context"
	"github.com/asoul-video/asoul-video/internal/db"
)

type Video struct{}

// NewVideoHandler creates a new Video router.
func NewVideoHandler() *Video {
	return &Video{}
}

func (*Video) List(ctx context.Context) {
	secUIDs := ctx.QueryStrings("secUID")
	keyword := ctx.Query("keyword")
	orderBy := "created_at"
	if ctx.Query("orderBy") != "" {
		orderBy = ctx.Query("orderBy")
	}
	order := ctx.Query("order")

	page := ctx.QueryInt("page")
	pageSize := ctx.QueryInt("pageSize")

	videos, err := db.Videos.List(ctx.Request().Context(), db.ListVideoOptions{
		SecUIDs:  secUIDs,
		Keyword:  keyword,
		OrderBy:  orderBy,
		Order:    order,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		log.Error("Failed to list video: %v", err)
		ctx.ServerError()
		return
	}

	ctx.Success(videos)
}

func (*Video) GetByID(ctx context.Context) {
	id := ctx.Param("id")
	video, err := db.Videos.GetByID(ctx.Request().Context(), id)
	if err != nil {
		log.Error("Failed to get video by ID: %v", err)
		ctx.ServerError()
		return
	}

	ctx.Success(video)
}

func (*Video) Random(ctx context.Context) {
	video, err := db.Videos.Random(ctx.Request().Context())
	if err != nil {
		log.Error("Failed to get video randomly: %v", err)
		ctx.ServerError()
		return
	}

	ctx.Success(video)
}
