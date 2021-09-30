// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"upper.io/db.v3/lib/sqlbuilder"

	"github.com/asoul-video/asoul-video/internal/dbutil"
)

var _ CommentsStore = (*comments)(nil)

var Comments CommentsStore

type CommentsStore interface {
	// Create creates a new comment with the given options.
	// If the comment already exists, it returns ErrCommentExists.
	Create(ctx context.Context, cid string, opts CreateCommentOptions) error
}

func NewCommentsStore(db sqlbuilder.Database) CommentsStore {
	return &comments{db}
}

type Comment struct {
	Cid           string          `db:"cid" json:"cid"`
	VideoID       string          `db:"video_id" json:"video_id"`
	DiggCount     int64           `db:"digg_count" json:"digg_count"`
	Text          string          `db:"text" json:"text"`
	TextClean     string          `db:"text_clean" json:"text_clean"`
	TextExtra     json.RawMessage `db:"text_extra" json:"text_extra"`
	UserNickname  string          `db:"user_nickname" json:"user_nickname"`
	UserAvatarURI string          `db:"user_avatar_uri" json:"user_avatar_uri"`
	UserSecUID    string          `db:"user_sec_uid" json:"user_sec_uid"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}

type comments struct {
	sqlbuilder.Database
}

type CreateCommentOptions struct {
	VideoID       string
	Text          string
	TextClean     string
	TextExtra     json.RawMessage
	UserNickname  string
	UserAvatarURI string
	UserSecUID    string
	CreatedAt     time.Time
}

var ErrCommentExists = errors.New("comment exists")

func (db *comments) Create(ctx context.Context, cid string, opts CreateCommentOptions) error {
	_, err := db.WithContext(ctx).InsertInto("comments").
		Columns("cid", "video_id", "text", "text_clean", "text_extra", "user_nickname", "user_avatar_uri", "user_sec_uid", "created_at").
		Values(cid, opts.VideoID, opts.Text, opts.TextClean, opts.TextExtra, opts.UserNickname, opts.UserAvatarURI, opts.UserSecUID, opts.CreatedAt).
		Exec()

	if err != nil {
		if dbutil.IsUniqueViolation(err, "comments_pkey") {
			return ErrCommentExists
		}
		return err
	}
	return nil
}
