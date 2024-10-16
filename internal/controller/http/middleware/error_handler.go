package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"

	"github.com/cewuandy/go-restful-dns/internal/domain"
)

type errorHandler struct {
}

func (h *errorHandler) HandleError(ctx *gin.Context) {
	ctx.Next()

	if len(ctx.Errors) == 0 {
		return
	}

	for _, err := range ctx.Errors {
		var e domain.Error
		_ = json.Unmarshal([]byte(err.Error()), &e)

		ctx.Header("Content-type", "application/problem+json")
		ctx.AbortWithStatusJSON(e.StatusCode, e)
	}

}

func NewErrorHandler(injector *do.Injector) (domain.ErrorHandler, error) {
	return &errorHandler{}, nil
}
