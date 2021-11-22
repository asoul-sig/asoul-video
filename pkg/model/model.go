// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"
)

type Statistic struct {
	Share   int64 `db:"share" json:"share"`
	Forward int64 `db:"forward" json:"forward"`
	Digg    int64 `db:"digg" json:"digg"`
	Play    int64 `db:"play" json:"play"`
	Comment int64 `db:"comment" json:"comment"`
}

type UpdateMember struct {
	SecUID    MemberSecUID `json:"sec_uid"`
	UID       string       `json:"uid"`
	UniqueID  string       `json:"unique_id"`
	ShortUID  string       `json:"short_uid"`
	Name      string       `json:"name"`
	AvatarURL string       `json:"avatar_url"`
	Signature string       `json:"signature"`
}

type CreateVideo struct {
	ID               string       `json:"id"`
	VID              string       `json:"vid"`
	AuthorSecUID     MemberSecUID `json:"author_sec_uid"`
	Description      string       `json:"description"`
	TextExtra        []string     `json:"text_extra"`
	OriginCoverURLs  []string     `json:"origin_cover_urls"`
	DynamicCoverURLs []string     `json:"dynamic_cover_urls"`
	IsDynamicCover   bool         `json:"is_dynamic_cover"`
	VideoHeight      int          `json:"video_height"`
	VideoWidth       int          `json:"video_width"`
	VideoDuration    int64        `json:"video_duration"`
	VideoRatio       string       `json:"video_ratio"`
	VideoURLs        []string     `json:"video_urls"`
	VideoCDNURL      string       `json:"video_cdn_url"`

	Statistic
}

type UpdateVideoMeta struct {
	ID               string    `json:"id"`
	VID              string    `json:"vid"`
	OriginCoverURLs  []string  `json:"origin_cover_urls"`
	DynamicCoverURLs []string  `json:"dynamic_cover_urls"`
	IsDynamicCover   bool      `json:"is_dynamic_cover"`
	CreatedAt        time.Time `json:"created_at"`

	Statistic
}

type CreateComment struct {
	Cid           string          `json:"cid"`
	VideoID       string          `json:"video_id"`
	Text          string          `json:"text"`
	TextClean     string          `json:"text_clean"`
	TextExtra     json.RawMessage `json:"text_extra"`
	UserNickname  string          `json:"user_nickname"`
	UserAvatarURI string          `json:"user_avatar_uri"`
	UserSecUID    string          `json:"user_sec_uid"`
	CreatedAt     time.Time       `json:"created_at"`
}

type UpdateFacePoint struct {
	ID          string          `json:"id"`
	FacePoints  json.RawMessage `json:"face_points"`
	CoverWidth  int             `json:"cover_width"`
	CoverHeight int             `json:"cover_height"`
}
