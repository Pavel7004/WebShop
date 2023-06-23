package v1

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/domain"
)

// CreateOrder godoc
// @Summary     Place new order
// @Description	Create order and record it
// @Tags        Orders
// @Accept		json
// @Produce     json
// @Param       req	  body  domain.CreateOrderRequest	true  "Request to create an order"
// @Success      200  {object}  string
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /shop/v1/orders/new [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendError(c, err)
		return
	}

	id, err := h.shop.CreateOrder(ctx, &req)
	if err != nil {
		h.SendError(c, err)
		return
	}

	span.SetTag("order_id", id)

	c.JSON(200, id)
}
