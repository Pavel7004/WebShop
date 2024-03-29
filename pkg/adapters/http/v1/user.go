package v1

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/domain"
)

// GetUser godoc
// @Summary      Get user
// @Description  Get user by ID
// @Tags         Users
// @Produce      json
// @Param        user_id  path  int  true  "user ID"
// @Success      200  {object}  domain.User
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /shop/v1/user/{user_id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	id := c.Param("user_id")

	span.SetTag("user_id", id)

	item, err := h.shop.GetUserById(ctx, id)
	if err != nil {
		h.SendError(c, err)
		return
	}

	c.JSON(200, item)
}

// RegisterUser godoc
// @Summary     Register user
// @Description	Register new user
// @Tags        Users
// @Accept		json
// @Produce     json
// @Param       req  body  domain.RegisterUserRequest	true  "Request to register new user"
// @Success      200  {object}  string
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /shop/v1/user/new [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	var req domain.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendError(c, err)
		return
	}

	id, err := h.shop.RegisterUser(ctx, &req)
	if err != nil {
		h.SendError(c, err)
		return
	}

	c.JSON(200, id)
}

// GetItemsByOwnerId godoc
// @Summary     Get items owned by 'user_id'
// @Description	Get all items that were created by user
// @Tags        Users
// @Produce     json
// @Success      200  {object}  []domain.Item
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /shop/v1/user/{user_id}/items [get]
func (h *Handler) GetItemsByOwnerId(c *gin.Context) {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	id := c.Param("user_id")

	span.SetTag("user_id", id)

	items, err := h.shop.GetItemsByOwnerId(ctx, id)
	if err != nil {
		h.SendError(c, err)
		return
	}

	c.JSON(200, items)
}

// GetRecentlyAddedUsers godoc
// @Summary     Get recenly added users
// @Description	Get last 2 added users
// @Tags        Users
// @Produce     json
// @Success      200  {object}  []domain.User
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /shop/v1/users/recent [get]
func (h *Handler) GetRecentlyAddedUsers(c *gin.Context) {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	users, err := h.shop.GetRecentlyAddedUsers(ctx, h.cfg.RecentUsersCount)
	if err != nil {
		h.SendError(c, err)
		return
	}

	c.JSON(200, users)
}
