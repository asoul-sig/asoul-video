// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/flamego/template"
)

type Landing struct{}

// NewLandingHandler creates a new Landing router.
func NewLandingHandler() *Landing {
	return &Landing{}
}

func (l *Landing) Home(t template.Template) {
	t.HTML(http.StatusOK, "index")
}
