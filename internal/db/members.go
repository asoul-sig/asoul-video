// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"github.com/pkg/errors"
	dbv3 "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"

	"github.com/asoul-video/asoul-video/pkg/model"
)

var _ MembersStore = (*members)(nil)

var Members MembersStore

type MembersStore interface {
	// Create creates a new member profile record with the given options
	// if the member's `name` `avatar_url` `signature` has been updated.
	Create(ctx context.Context, opts CreateMemberOptions) error
	// GetBySecID returns the latest member profile with the given SecUID.
	GetBySecID(ctx context.Context, secUID model.MemberSecUID) (*Member, error)
	// GetBySecIDs returns the members' profile with the given SecUIDs.
	// It will be ignored if the member does not exist.
	GetBySecIDs(ctx context.Context, secUIDs ...model.MemberSecUID) ([]*Member, error)
}

func NewMembersStore(db sqlbuilder.Database) MembersStore {
	return &members{db}
}

type Member struct {
	SecUID    model.MemberSecUID `db:"sec_uid"`
	UID       string             `db:"uid"`
	UniqueID  string             `db:"unique_id"`
	ShortUID  string             `db:"short_id"`
	Name      string             `db:"name"`
	AvatarURL string             `db:"avatar_url"`
	Signature string             `db:"signature"`
}

type members struct {
	sqlbuilder.Database
}

type CreateMemberOptions struct {
	SecUID    model.MemberSecUID
	UID       string
	UniqueID  string
	ShortUID  string
	Name      string
	AvatarURL string
	Signature string
}

func (db *members) Create(ctx context.Context, opts CreateMemberOptions) error {
	member, err := db.GetBySecID(ctx, opts.SecUID)
	if err == nil {
		opts.SecUID = member.SecUID
		opts.UID = member.UID
		opts.UniqueID = member.UniqueID
		opts.ShortUID = member.ShortUID

		if opts.Name == member.Name && opts.AvatarURL == member.AvatarURL && opts.Signature == member.Signature {
			return nil
		}
	}

	_, err = db.WithContext(ctx).InsertInto("members").
		Columns("uid", "unique_id", "sec_uid", "short_id", "name", "avatar_url", "signature").
		Values(opts.UID, opts.UniqueID, opts.SecUID, opts.ShortUID, opts.Name, opts.AvatarURL, opts.Signature).
		Exec()
	return err
}

var ErrMemberNotFound = errors.New("member dose not exist")

func (db *members) GetBySecID(ctx context.Context, secUID model.MemberSecUID) (*Member, error) {
	var member Member

	err := db.WithContext(ctx).SelectFrom("members").
		Where("sec_uid = ?", secUID).
		OrderBy("created_at DESC").
		One(&member)
	if err != nil {
		if err == dbv3.ErrNoMoreRows {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	return &member, nil
}

func (db *members) GetBySecIDs(ctx context.Context, secUIDs ...model.MemberSecUID) ([]*Member, error) {
	var members []*Member

	return members, db.WithContext(ctx).SelectFrom("members").
		Where("sec_uid IN ?", secUIDs).
		OrderBy("created_at DESC").
		All(&members)
}
