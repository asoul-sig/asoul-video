// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"time"

	"upper.io/db.v3/lib/sqlbuilder"
)

var _ StatisticsStore = (*statistics)(nil)

var Statistics StatisticsStore

type StatisticsStore interface {
	// Create creates a new statistic record with the given options.
	Create(ctx context.Context, id string, opts CreateStatisticOptions) error
}

func NewStatisticsStore(db sqlbuilder.Database) StatisticsStore {
	return &statistics{db}
}

type Statistic struct {
	ID        string    `db:"id" json:"id"`
	Share     int64     `db:"share" json:"share"`
	Forward   int64     `db:"forward" json:"forward"`
	Digg      int64     `db:"digg" json:"digg"`
	Play      int64     `db:"play" json:"play"`
	Comment   int64     `db:"comment" json:"comment"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

type statistics struct {
	sqlbuilder.Database
}

type CreateStatisticOptions struct {
	Share   int64
	Forward int64
	Digg    int64
	Play    int64
	Comment int64
}

func (db *statistics) Create(ctx context.Context, id string, opts CreateStatisticOptions) error {
	if opts.Share+opts.Forward+opts.Digg+opts.Play+opts.Comment == 0 {
		return nil
	}

	var statistic Statistic
	if err := db.WithContext(ctx).SelectFrom("statistic").Where(
		"id = ? AND share = ? AND forward = ? AND digg = ? AND play = ? AND comment = ?",
		id, opts.Share, opts.Forward, opts.Digg, opts.Play, opts.Comment).One(&statistic); err == nil {
		return nil
	}

	_, err := db.WithContext(ctx).InsertInto("statistics").
		Columns("id", "share", "forward", "digg", "play", "comment").
		Values(id, opts.Share, opts.Forward, opts.Digg, opts.Play, opts.Comment).
		Exec()
	return err
}
