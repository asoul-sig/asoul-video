// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"upper.io/db.v3/lib/sqlbuilder"

	"github.com/asoul-video/asoul-video/internal/dbutil"
)

type VideoStatus string

var (
	VideoStatusAvailable   VideoStatus = "available"
	VideoStatusUnavailable VideoStatus = "unavailable"
)

var _ VideoURLsStore = (*videoURLs)(nil)

var VideoURLs VideoURLsStore

type VideoURLsStore interface {
	// Create creates a new video url with the given video ID and url.
	Create(ctx context.Context, videoID, url string) error
	// GetByVideoID returns the available video urls with the given video ID.
	GetByVideoID(ctx context.Context, videoID string) ([]string, error)
	// SetStatus set the video url status.
	SetStatus(ctx context.Context, videoURL string, status VideoStatus) error
}

func NewVideoURLsStore(db sqlbuilder.Database) VideoURLsStore {
	return &videoURLs{db}
}

type VideoURL struct {
	VideoID     string    `db:"video_id"`
	URL         string    `db:"url"`
	Status      string    `db:"status"`
	LastCheckAt time.Time `db:"last_check_at"`
}

type videoURLs struct {
	sqlbuilder.Database
}

func (db *videoURLs) Create(ctx context.Context, videoID, u string) error {
	videoURL, err := url.Parse(u)
	if err != nil {
		return errors.Wrap(err, "parse URL")
	}
	if strings.HasSuffix(videoURL.Host, "douyinvod.com") { // Remove temporary video URL.
		return nil
	}

	_, err = db.WithContext(ctx).InsertInto("video_urls").
		Columns("video_id", "url", "status").
		Values(videoID, u, VideoStatusAvailable).
		Exec()
	if err != nil {
		if dbutil.IsUniqueViolation(err, "video_urls_pkey") {
			return ErrVideoExists
		}
		return err
	}
	return nil
}

func (db *videoURLs) GetByVideoID(ctx context.Context, videoID string) ([]string, error) {
	var urls []string
	return urls, db.WithContext(ctx).Select("url").From("video_urls").
		Where("video_id = ? AND status = ?", videoID, VideoStatusAvailable).
		All(&urls)
}

func (db *videoURLs) SetStatus(ctx context.Context, videoURL string, status VideoStatus) error {
	q := db.WithContext(ctx).Update("video_urls")

	switch status {
	case VideoStatusAvailable:
		q = q.Set("status", VideoStatusAvailable, "last_check_at", time.Now())
	case VideoStatusUnavailable:
		q = q.Set("status", VideoStatusUnavailable, "last_check_at", time.Now())
	default:
		return errors.Errorf("unexpected status: %q", status)
	}

	_, err := q.Where("video_id = ?", videoURL).Exec()
	if err != nil {
		return errors.Wrap(err, "update status")
	}
	return nil
}
