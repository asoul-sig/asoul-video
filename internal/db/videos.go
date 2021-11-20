// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	dbv3 "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"

	"github.com/asoul-video/asoul-video/internal/dbutil"
	"github.com/asoul-video/asoul-video/pkg/model"
)

var _ VideosStore = (*videos)(nil)

var Videos VideosStore

type VideosStore interface {
	// Create creates a new video record with the given options.
	Create(ctx context.Context, id string, opts CreateVideoOptions) error
	// Update updates the video with the given options.
	Update(ctx context.Context, id string, opts UpdateVideoOptions) error
	// GetByID returns video with the given id, it returns `ErrVideoNotFound` error if video does not exist.
	GetByID(ctx context.Context, id string) (*Video, error)
	// List returns the video list.
	List(ctx context.Context, opts ListVideoOptions) ([]*Video, error)
	// ListIDs returns all the video IDs.
	ListIDs(ctx context.Context) ([]string, error)
	// Random returns a video randomly.
	Random(ctx context.Context) (*Video, error)
}

func NewVideosStore(db sqlbuilder.Database) VideosStore {
	return &videos{db}
}

type Video struct {
	ID               string             `db:"id" json:"id"`
	VID              string             `db:"vid" json:"vid"`
	AuthorSecUID     model.MemberSecUID `db:"author_sec_id" json:"author_sec_uid"`
	Author           *Member            `db:",inline" json:"author"`
	Statistic        model.Statistic    `db:",inline" json:"statistic"`
	Description      string             `db:"description" json:"description"`
	TextExtra        []string           `db:"text_extra" json:"text_extra"`
	OriginCoverURLs  []string           `db:"origin_cover_urls" json:"origin_cover_urls"`
	DynamicCoverURLs []string           `db:"dynamic_cover_urls" json:"dynamic_cover_urls"`
	IsDynamicCover   bool               `db:"is_dynamic_cover" json:"is_dynamic_cover"`
	VideoHeight      int                `db:"video_height" json:"video_height"`
	VideoWidth       int                `db:"video_width" json:"video_width"`
	VideoDuration    int64              `db:"video_duration" json:"video_duration"`
	VideoRatio       string             `db:"video_ratio" json:"video_ratio"`
	CreatedAt        time.Time          `db:"created_at" json:"created_at"`
}

type videos struct {
	sqlbuilder.Database
}

type CreateVideoOptions struct {
	VID              string
	AuthorSecUID     model.MemberSecUID
	Description      string
	TextExtra        []string
	OriginCoverURLs  []string
	DynamicCoverURLs []string
	IsDynamicCover   bool
	VideoHeight      int
	VideoWidth       int
	VideoDuration    int64
	VideoRatio       string
}

var ErrVideoExists = errors.New("duplicate video")

func (db *videos) Create(ctx context.Context, id string, opts CreateVideoOptions) error {
	_, err := db.WithContext(ctx).InsertInto("videos").
		Columns("id", "vid", "author_sec_id", "description", "text_extra", "origin_cover_urls", "dynamic_cover_urls", "is_dynamic_cover", "video_height", "video_width", "video_duration", "video_ratio").
		Values(id, opts.VID, opts.AuthorSecUID, opts.Description, opts.TextExtra, opts.OriginCoverURLs, opts.DynamicCoverURLs, opts.IsDynamicCover, opts.VideoHeight, opts.VideoWidth, opts.VideoDuration, opts.VideoRatio).
		Exec()
	if err != nil {
		if dbutil.IsUniqueViolation(err, "videos_pkey") {
			if err := db.Update(ctx, id, UpdateVideoOptions{
				VID:              opts.VID,
				IsDynamicCover:   opts.IsDynamicCover,
				OriginCoverURLs:  opts.OriginCoverURLs,
				DynamicCoverURLs: opts.DynamicCoverURLs,
			}); err != nil {
				return errors.Wrap(err, "update video")
			}
		} else {
			return errors.Wrap(err, "create video")
		}
	}
	return nil
}

type UpdateVideoOptions struct {
	VID              string
	IsDynamicCover   bool
	OriginCoverURLs  []string
	DynamicCoverURLs []string
	CreatedAt        time.Time
}

func (db *videos) Update(ctx context.Context, id string, opts UpdateVideoOptions) error {
	_, err := db.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get video by id")
	}

	updateSets := make([]interface{}, 0, 10)
	updateSets = append(updateSets, "is_dynamic_cover", opts.IsDynamicCover)

	if opts.VID != "" {
		updateSets = append(updateSets, "vid", opts.VID)
	}
	if len(opts.OriginCoverURLs) > 0 {
		updateSets = append(updateSets, "origin_cover_urls", opts.OriginCoverURLs)
	}
	if len(opts.DynamicCoverURLs) > 0 {
		updateSets = append(updateSets, "dynamic_cover_urls", opts.DynamicCoverURLs)
	}
	if !opts.CreatedAt.IsZero() {
		updateSets = append(updateSets, "created_at", opts.CreatedAt)
	}
	if len(updateSets) == 0 {
		return nil
	}

	updateSets = append(updateSets, "updated_at", time.Now())
	_, err = db.WithContext(ctx).Update("videos").
		Set(updateSets...).
		Where("id = ?", id).Exec()
	return err
}

var ErrVideoNotFound = errors.New("video dose not exist")

func (db *videos) GetByID(ctx context.Context, id string) (*Video, error) {
	var video Video
	if err := db.WithContext(ctx).SelectFrom("video_list").Where("id = ?", id).One(&video); err != nil {
		if err == dbv3.ErrNoMoreRows {
			return nil, ErrVideoNotFound
		}
		return nil, err
	}
	return &video, nil
}

type ListVideoOptions struct {
	Keyword string
	SecUIDs []string
	OrderBy string
	Order   string

	Page     int
	PageSize int
}

func (db *videos) List(ctx context.Context, opts ListVideoOptions) ([]*Video, error) {
	if opts.OrderBy != "video_duration" && opts.OrderBy != "created_at" {
		opts.OrderBy = ""
	}

	if opts.Order != "asc" {
		opts.Order = "desc"
	}

	if opts.Page <= 0 {
		opts.Page = 1
	}

	if opts.PageSize <= 0 || opts.PageSize >= 30 {
		opts.PageSize = 30
	}

	query := db.WithContext(ctx).SelectFrom("video_list")

	if len(opts.SecUIDs) != 0 {
		query = query.Where("author_sec_id IN ?", opts.SecUIDs)
	}

	if opts.Keyword != "" {
		query = query.And("description ILIKE ?", "%"+opts.Keyword+"%")
	}

	if opts.OrderBy != "" {
		query = query.OrderBy(opts.OrderBy + " " + opts.Order)
	}

	// Pagination
	query = query.Limit(opts.PageSize).Offset((opts.Page - 1) * opts.PageSize)

	var videos []*Video
	if err := query.All(&videos); err != nil {
		return nil, errors.Wrap(err, "get videos")
	}
	return videos, nil
}

func (db *videos) ListIDs(ctx context.Context) ([]string, error) {
	var idRows []struct {
		ID string `db:"id"`
	}
	if err := db.WithContext(ctx).Select("id").From("videos").OrderBy("created_at").All(&idRows); err != nil {
		return nil, errors.Wrap(err, "select")
	}

	ids := make([]string, 0, len(idRows))
	for _, row := range idRows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

func (db *videos) Random(ctx context.Context) (*Video, error) {
	var count int
	row, err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM video_list;")
	if err != nil {
		return nil, errors.Wrap(err, "count")
	}
	if err := row.Scan(&count); err != nil {
		return nil, errors.Wrap(err, "count")
	}

	var video Video
	if err := db.SelectFrom("video_list").Offset(rand.Intn(count)).Limit(1).One(&video); err != nil {
		return nil, errors.Wrap(err, "get video")
	}
	return &video, nil
}
