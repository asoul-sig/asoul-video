// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"math/rand"
	"net/url"
	"strings"
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
	// GetByID returns video with the given id, it returns `ErrVideoNotFound` error if video does not exist.
	GetByID(ctx context.Context, id string) (*Video, error)
	// List returns the video list.
	List(ctx context.Context, opts ListVideoOptions) ([]*Video, error)
	// Random returns a video randomly.
	Random(ctx context.Context) (*Video, error)
}

func NewVideosStore(db sqlbuilder.Database) VideosStore {
	return &videos{db}
}

type Video struct {
	ID               string             `db:"id" json:"id"`
	AuthorSecUID     model.MemberSecUID `db:"author_sec_id" json:"author_sec_uid"`
	Author           *Member            `db:"-" json:"author"`
	Description      string             `db:"description" json:"description"`
	TextExtra        []string           `db:"text_extra" json:"text_extra"`
	OriginCoverURLs  []string           `db:"origin_cover_urls" json:"origin_cover_urls"`
	DynamicCoverURLs []string           `db:"dynamic_cover_urls" json:"dynamic_cover_urls"`
	VideoHeight      int                `db:"video_height" json:"video_height"`
	VideoWidth       int                `db:"video_width" json:"video_width"`
	VideoDuration    int64              `db:"video_duration" json:"video_duration"`
	VideoRatio       string             `db:"video_ratio" json:"video_ratio"`
	VideoURLs        []string           `db:"video_urls" json:"video_urls"`
	VideoCDNURL      string             `db:"video_cdn_url" json:"-"`
	CreatedAt        time.Time          `db:"created_at" json:"created_at"`
}

type videos struct {
	sqlbuilder.Database
}

type CreateVideoOptions struct {
	AuthorSecUID     model.MemberSecUID
	Description      string
	TextExtra        []string
	OriginCoverURLs  []string
	DynamicCoverURLs []string
	VideoHeight      int
	VideoWidth       int
	VideoDuration    int64
	VideoRatio       string
	VideoURLs        []string
	VideoCDNURL      string
	CreatedAt        time.Time
}

var ErrVideoExists = errors.New("duplicate video")

func (db *videos) Create(ctx context.Context, id string, opts CreateVideoOptions) error {
	if opts.CreatedAt.IsZero() {
		opts.CreatedAt = time.Now()
	}

	videoURLs := make([]string, 0, len(opts.VideoURLs))
	for _, u := range videoURLs {
		videoURL, err := url.Parse(u)
		if err != nil {
			continue
		}
		if strings.HasSuffix(videoURL.Host, "douyinvod.com") { // Remove temporary video URL.
			continue
		}

		videoURLs = append(videoURLs, u)
	}

	_, err := db.WithContext(ctx).InsertInto("videos").
		Columns("id", "author_sec_id", "description", "text_extra", "origin_cover_urls", "dynamic_cover_urls", "video_height", "video_width", "video_duration", "video_ratio", "video_urls", "video_cdn_url", "created_at").
		Values(id, opts.AuthorSecUID, opts.Description, opts.TextExtra, opts.OriginCoverURLs, opts.DynamicCoverURLs, opts.VideoHeight, opts.VideoWidth, opts.VideoDuration, opts.VideoRatio, videoURLs, opts.VideoCDNURL, opts.CreatedAt).
		Exec()
	if err != nil {
		if dbutil.IsUniqueViolation(err, "videos_pkey") {
			return ErrVideoExists
		}
		return err
	}

	return nil
}

var ErrVideoNotFound = errors.New("video dose not exist")

func (db *videos) GetByID(ctx context.Context, id string) (*Video, error) {
	var video Video

	err := db.WithContext(ctx).SelectFrom("videos").
		Where("id = ?", id).
		One(&video)
	if err != nil {
		if err == dbv3.ErrNoMoreRows {
			return nil, ErrVideoNotFound
		}
		return nil, err
	}

	memberStore := NewMembersStore(db)
	member, err := memberStore.GetBySecID(ctx, video.AuthorSecUID)
	if err != nil {
		return nil, errors.Wrap(err, "get member by sec id")
	}
	video.Author = member

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

	query := db.WithContext(ctx).SelectFrom("videos")

	if len(opts.SecUIDs) != 0 {
		query = query.Where("author_sec_id IN ?", opts.SecUIDs)
	}

	if opts.Keyword != "" {
		query = query.And("description ILIKE ?", opts.Keyword)
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

	memberSecIDSets := make(map[model.MemberSecUID]struct{})
	for _, video := range videos {
		memberSecIDSets[video.AuthorSecUID] = struct{}{}
	}

	memberSecIDs := make([]model.MemberSecUID, 0, len(memberSecIDSets))
	for memberSecID := range memberSecIDSets {
		memberSecIDs = append(memberSecIDs, memberSecID)
	}

	memberStore := NewMembersStore(db)
	members, err := memberStore.GetBySecIDs(ctx, memberSecIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "get videos' members")
	}

	// FIXME: maybe we can just query it from database with the videos?
	memberSets := make(map[model.MemberSecUID]*Member, len(members))
	for _, member := range members {
		member := member
		memberSets[member.SecUID] = member
	}
	for _, video := range videos {
		video.Author = memberSets[video.AuthorSecUID]
	}

	return videos, nil
}

func (db *videos) Random(ctx context.Context) (*Video, error) {
	var count int
	row, err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM videos;")
	if err != nil {
		return nil, errors.Wrap(err, "count")
	}
	if err := row.Scan(&count); err != nil {
		return nil, errors.Wrap(err, "count")
	}

	var video Video
	if err := db.SelectFrom("videos").Offset(rand.Intn(count)).Limit(1).One(&video); err != nil {
		return nil, errors.Wrap(err, "get video")
	}

	memberStore := NewMembersStore(db)
	member, err := memberStore.GetBySecID(ctx, video.AuthorSecUID)
	if err != nil {
		return nil, errors.Wrap(err, "get member by re")
	}
	video.Author = member

	return &video, nil
}
