// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"net/http"

	"github.com/flamego/flamego"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

// Context represents context of a request.
type Context struct {
	flamego.Context
}

func (c *Context) Success(data interface{}) {
	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)

	err := jsoniter.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"data": data,
		},
	)
	if err != nil {
		log.Error("Failed to encode: %v", err)
	}
}

func (c *Context) ServerError() {
	c.Error(http.StatusInternalServerError, errors.New("internal server error"))
}

func (c *Context) Error(statusCode int, error error) {
	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(statusCode)

	err := jsoniter.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"msg": error.Error(),
		},
	)
	if err != nil {
		log.Error("Failed to encode: %v", err)
	}
}

// Contexter initializes a classic context for a request.
func Contexter() flamego.Handler {
	return func(ctx flamego.Context) {
		c := Context{
			Context: ctx,
		}
		c.Map(c)
	}
}
