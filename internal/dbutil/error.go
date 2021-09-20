// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
)

func isConstraintViolation(err error, code, constraint string) bool {
	pgErr, ok := errors.Cause(err).(*pgconn.PgError)
	if !ok {
		return false
	}
	return pgErr.Code == code && pgErr.ConstraintName == constraint
}

func IsUniqueViolation(err error, constraint string) bool {
	return isConstraintViolation(err, "23505", constraint)
}
