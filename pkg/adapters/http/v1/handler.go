package v1

import (
	"github.com/Pavel7004/WebShop/pkg/components"
	"github.com/Pavel7004/WebShop/pkg/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	shop components.Shop
}

func New(shop components.Shop) *Handler {
	return &Handler{
		shop: shop,
	}
}

func (h *Handler) SendError(c *gin.Context, err error) {
	if e, ok := err.(*domain.Error); ok {
		c.JSON(e.CodeHTTP, e)
	} else {
		c.JSON(500, &domain.Error{
			Code:    "unknown_error",
			Message: err.Error(),
		})
	}
}
