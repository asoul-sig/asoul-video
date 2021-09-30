// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"os"

	"github.com/pkg/errors"
	"upper.io/db.v3/postgresql"

	"github.com/asoul-video/asoul-video/internal/dbutil"
	"github.com/asoul-video/asoul-video/migrations"
)

// Init connects to the database and migrate tables.
func Init() error {
	db, err := postgresql.Open(postgresql.ConnectionURL{
		Host:     os.Getenv("PGHOST"),
		User:     os.Getenv("PGUSER"),
		Password: os.Getenv("PGPASSWORD"),
		Database: os.Getenv("PGDATABASE"),
		Options: map[string]string{
			"sslmode": os.Getenv("PGSSLMODE"),
		},
	})
	if err != nil {
		return errors.Wrap(err, "open database")
	}

	_, err = dbutil.Migrate(db.Driver().(*sql.DB), migrations.Migrations)
	if err != nil {
		return errors.Wrap(err, "migrate")
	}

	Members = NewMembersStore(db)
	Videos = NewVideosStore(db)
	VideoURLs = NewVideoURLsStore(db)
	Statistics = NewStatisticsStore(db)
	Comments = NewCommentsStore(db)

	return nil
}
