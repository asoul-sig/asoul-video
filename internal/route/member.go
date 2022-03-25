// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	log "unknwon.dev/clog/v2"

	"github.com/asoul-sig/asoul-video/internal/context"
	"github.com/asoul-sig/asoul-video/internal/db"
	"github.com/asoul-sig/asoul-video/pkg/model"
)

type Member struct{}

// NewMemberHandler creates a new Member router.
func NewMemberHandler() *Member {
	return &Member{}
}

func (*Member) List(ctx context.Context) {
	members, err := db.Members.List(ctx.Request().Context())
	if err != nil {
		log.Error("Failed to get member list: %v", err)
		ctx.ServerError()
		return
	}

	ctx.Success(members)
}

func (*Member) GetBySecUID(ctx context.Context) {
	secUID := ctx.Param("secUID")

	member, err := db.Members.GetBySecID(ctx.Request().Context(), model.MemberSecUID(secUID))
	if err != nil {
		if err == db.ErrMemberNotFound {
			ctx.Error(http.StatusNotFound, err)
			return
		}

		log.Error("Failed to get member by sec_uid: %v", err)
		ctx.ServerError()
		return
	}

	ctx.Success(member)
}
