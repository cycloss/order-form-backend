package main

import (
	"log"

	"github.com/cycloss/aj-bell-test/order-api/server"
	"github.com/cycloss/aj-bell-test/share"
	"github.com/gin-gonic/gin"
)

const serverBindAddress = "0.0.0.0:80"

func main() {

	db := share.MustConnectDb()
	hr := server.NewHandleWrapper(db)

	router := gin.Default()
	orderApiGroup := router.Group("/order-api")
	orderApiGroup.POST("/buy", hr.BuyHandler)
	orderApiGroup.POST("/invest", nil)
	orderApiGroup.POST("/sell", nil)
	orderApiGroup.POST("/raise", nil)
	orderApiGroup.GET("/asset-holdings", nil)
	orderApiGroup.GET("/currency-holdings", nil)

	log.Printf("Starting server on %s", serverBindAddress)
	router.Run(serverBindAddress)
}

func Bar() {

}
