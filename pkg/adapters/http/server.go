package http

import (
	v1 "github.com/Pavel7004/WebShop/pkg/adapters/http/v1"
	"github.com/Pavel7004/WebShop/pkg/components"
	"github.com/Pavel7004/WebShop/pkg/infra/config"
	"github.com/gin-gonic/gin"
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
		v1.GET("/item/:item_id", s.v1.GetItem)              // -
		v1.POST("/item/new", s.v1.AddItem)                  // -
		v1.GET("/items", s.v1.GetItems)                     // -
		v1.GET("/items/recent", s.v1.GetRecentlyAddedItems) // -

		v1.GET("/user/:user_id", s.v1.GetUser) // -
		v1.GET("/user/new", s.v1.RegisterUser) // -
	}

	// query ?a=1&b=2 <- GET, DELETE не имеют тела
	// body json {"a": 1, "b": 2} POST, PUT - инициализация больших объектов
	// {last_name: ..., surname: ..., balance: ..., date: ...}
	// body x-www-form-urlencoded (a=1&b=2 in body) POST, PUT
	// path /item/{item_id}
}
