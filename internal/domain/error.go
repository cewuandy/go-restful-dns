package domain

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Err        error  `json:"-"`
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type ErrorHandler interface {
	HandleError(ctx *gin.Context)
}
