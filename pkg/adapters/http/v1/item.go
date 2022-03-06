package v1

import (
	"context"
	"strconv"
	"time"

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
	id := c.Param("item_id")

	item, err := h.shop.GetItemById(context.Background(), id)
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
	var req domain.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendError(c, err)
		return
	}

	id, err := h.shop.AddItem(context.Background(), &req)
	if err != nil {
		h.SendError(c, err)
		return
	}

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
	var (
		fromStr = c.DefaultQuery("from", "0")
		toStr   = c.DefaultQuery("to", "10")
	)

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

	items, err := h.shop.GetItemsByPrice(context.Background(), from, to)
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
	periodStr := c.DefaultQuery("period", time.Now().Sub(time.Now().AddDate(0, 0, -3)).String())

	period, err := time.ParseDuration(periodStr)
	if err != nil {
		h.SendError(c, err)
		return
	}

	items, err := h.shop.GetRecentlyAddedItems(context.Background(), period)
	if err != nil {
		h.SendError(c, err)
		return
	}

	c.JSON(200, items)
}
