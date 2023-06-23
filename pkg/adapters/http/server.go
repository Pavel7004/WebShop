package http

// @title           WebShop API
// @version         0.1
// @description     This is an API for online store

// @contact.name   Kovalev Pavel
// @contact.email  kovalev5690@gmail.com

// @license.name   GPL-3.0
// @license.url    https://www.gnu.org/licenses/gpl-3.0.html

// @host      localhost:8080

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/Pavel7004/WebShop/pkg/adapters/http/v1"
	"github.com/Pavel7004/WebShop/pkg/components"
	"github.com/Pavel7004/WebShop/pkg/infra/config"
)

type Server struct {
	router    *gin.Engine
	isRunning bool

	v1 *v1.Handler
}

func New(shop components.Shop, cfg *config.Config) *Server {
	server := new(Server)

	server.router = gin.New()
	server.isRunning = false
	server.v1 = v1.New(shop, cfg)

	server.prepareRouter()

	return server
}

func (s *Server) Run() error {
	s.isRunning = true
	return s.router.Run(":8080")
}

func (s *Server) prepareRouter() {
	v1 := s.router.Group("/shop/v1")
	{
		v1.GET("/items/:item_id", s.v1.GetItem)             // -
		v1.POST("/items/new", s.v1.AddItem)                 // -
		v1.PUT("/items/:item_id", s.v1.UpdateItem)          // -
		v1.GET("/items", s.v1.GetItems)                     // -
		v1.GET("/items/recent", s.v1.GetRecentlyAddedItems) // -

		v1.GET("/user/:user_id", s.v1.GetUser)                 // -
		v1.POST("/user/new", s.v1.RegisterUser)                // -
		v1.GET("/user/:user_id/items", s.v1.GetItemsByOwnerId) // -
		v1.GET("/users/recent", s.v1.GetRecentlyAddedUsers)    // -

		v1.POST("/orders/new", s.v1.CreateOrder) // -
	}

	// query ?a=1&b=2 <- GET, DELETE не имеют тела
	// body json {"a": 1, "b": 2} POST, PUT - инициализация больших объектов
	// {last_name: ..., surname: ..., balance: ..., date: ...}
	// body x-www-form-urlencoded (a=1&b=2 in body) POST, PUT
	// path /item/{item_id}
	// browser saves GET queries to history; POST, PUT aren't saved
}
