// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"github.com/pkg/errors"
	dbv3 "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"

	"github.com/asoul-sig/asoul-video/pkg/model"
)

var _ MembersStore = (*members)(nil)

var Members MembersStore

type MembersStore interface {
	// Upsert creates a new member profile record with the given options,
	// it updates the `name` `avatar_url` `signature` field if the member is exists.
	Upsert(ctx context.Context, opts UpsertMemberOptions) error
	// GetBySecID returns the latest member profile with the given SecUID.
	GetBySecID(ctx context.Context, secUID model.MemberSecUID) (*Member, error)
	// GetBySecIDs returns the members' profile with the given SecUIDs.
	// It will be ignored if the member does not exist.
	GetBySecIDs(ctx context.Context, secUIDs ...model.MemberSecUID) ([]*Member, error)
	// List returns all the members.
	List(ctx context.Context) ([]*Member, error)
}

func NewMembersStore(db sqlbuilder.Database) MembersStore {
	return &members{db}
}

type Member struct {
	SecUID    model.MemberSecUID `db:"sec_uid" json:"sec_uid"`
	UID       string             `db:"uid" json:"uid"`
	UniqueID  string             `db:"unique_id" json:"unique_id"`
	ShortUID  string             `db:"short_id" json:"short_uid"`
	Name      string             `db:"name" json:"name"`
	AvatarURL string             `db:"avatar_url" json:"avatar_url"`
	Signature string             `db:"signature" json:"signature"`
}

type members struct {
	sqlbuilder.Database
}

type UpsertMemberOptions struct {
	SecUID    model.MemberSecUID
	UID       string
	UniqueID  string
	ShortUID  string
	Name      string
	AvatarURL string
	Signature string
}

func (db *members) Upsert(ctx context.Context, opts UpsertMemberOptions) error {
	_, err := db.GetBySecID(ctx, opts.SecUID)
	if err == nil {
		if _, err := db.WithContext(ctx).Update("members").Set(
			"name", opts.Name,
			"avatar_url", opts.AvatarURL,
			"signature", opts.Signature).
			Where("sec_uid = ?", opts.SecUID).Exec(); err != nil {
			return errors.Wrap(err, "update member")
		}
		return nil

	} else if err != ErrMemberNotFound {
		return errors.Wrap(err, "get member by sec_id")
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

func (db *members) List(ctx context.Context) ([]*Member, error) {
	var members []*Member
	return members, db.WithContext(ctx).SelectFrom("members").OrderBy("uid").All(&members)
}
