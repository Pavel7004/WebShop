package v1

import (
	"context"
	"strconv"
	"time"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/domain"
	"github.com/gin-gonic/gin"
)

// GetItem godoc
// @Summary      Get item
// @Description  Get item by ID
// @Tags         Items
// @Produce      json
// @Param        item_id   path      int  true  "Item ID"
// @Success      200  {object}  domain.Item
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /shop/v1/item/{item_id} [get]
func (h *Handler) GetItem(c *gin.Context) {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	id := c.Param("item_id")

	span.SetTag("item_id", id)

	item, err := h.shop.GetItemById(ctx, id)
	if err != nil {
		h.SendError(c, err)
		return
	}

	c.JSON(200, item)
}

// AddItem godoc
// @Summary     Add item
// @Description	Add item
// @Tags        Items
// @Accept		json
// @Produce     json
// @Param       req	  body  domain.AddItemRequest	true  "Request to add an item"
// @Success      200  {object}  string
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /shop/v1/item [post]
func (h *Handler) AddItem(c *gin.Context) {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	var req domain.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendError(c, err)
		return
	}

	span.SetTag("item_request", req)

	id, err := h.shop.AddItem(ctx, &req)
	if err != nil {
		h.SendError(c, err)
		return
	}

	span.SetTag("item_id", id)

	c.JSON(200, id)
}

// GetItemsByPrice godoc
// @Summary     Get items within price range
// @Description	Get items with specified price range
// @Tags        Items
// @Produce     json
// @Param       from	query	float64	false  "Price lower bound"
// @Param       to		query	float64	false  "Price upper bound"
// @Success      200  {object}  []domain.Item
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /shop/v1/items [get]
func (h *Handler) GetItems(c *gin.Context) {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	var (
		fromStr = c.DefaultQuery("from", "0")
		toStr   = c.DefaultQuery("to", "99999999")
	)

	span.SetTag("from_query", fromStr)
	span.SetTag("to_query", toStr)

	from, err := strconv.ParseFloat(fromStr, 64)
	if err != nil {
		h.SendError(c, err)
		return
	}

	to, err := strconv.ParseFloat(toStr, 64)
	if err != nil {
		h.SendError(c, err)
		return
	}

	span.SetTag("from_parsed", from)
	span.SetTag("to_parsed", to)

	items, err := h.shop.GetItemsByPrice(ctx, from, to)
	if err != nil {
		h.SendError(c, err)
		return
	}

	c.JSON(200, items)
}

// GetRecentlyAddedItems godoc
// @Summary     Get recenly added items
// @Description	Get items that was added within last 3 days
// @Tags        Items
// @Produce     json
// @Param       period	query	time.Duration	false  "Price lower bound"
// @Success      200  {object}  []domain.Item
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /shop/v1/items/recent [get]
func (h *Handler) GetRecentlyAddedItems(c *gin.Context) {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	periodStr := c.DefaultQuery("period", (72 * time.Hour).String())

	span.SetTag("period_query", periodStr)

	period, err := time.ParseDuration(periodStr)
	if err != nil {
		h.SendError(c, err)
		return
	}

	span.SetTag("period_parsed", period)

	items, err := h.shop.GetRecentlyAddedItems(ctx, period)
	if err != nil {
		h.SendError(c, err)
		return
	}

	c.JSON(200, items)
}
