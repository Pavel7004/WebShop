package main

import (
	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo"
	"github.com/Pavel7004/WebShop/pkg/adapters/http"
	"github.com/Pavel7004/WebShop/pkg/components/shop"
)

func main() {
	db := mongo.New()

	shop := shop.New(db)
	server := http.New(shop)

	server.Run()
}
