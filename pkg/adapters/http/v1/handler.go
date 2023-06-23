package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/Pavel7004/WebShop/pkg/components"
	"github.com/Pavel7004/WebShop/pkg/domain"
	"github.com/Pavel7004/WebShop/pkg/infra/config"
)

type Handler struct {
	shop components.Shop
	cfg  *config.Config
}

func New(shop components.Shop, cfg *config.Config) *Handler {
	return &Handler{
		shop: shop,
		cfg:  cfg,
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
