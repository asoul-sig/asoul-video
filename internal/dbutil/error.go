// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"github.com/Shyp/go-dberror"
)

func IsUniqueViolation(err error, constraint string) bool {
	switch e := dberror.GetError(err).(type) {
	case *dberror.Error:
		return e.Constraint == constraint
	default:
		return false
	}
}
