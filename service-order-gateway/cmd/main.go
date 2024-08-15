package main

import (
	"github.com/gin-gonic/gin"
	"service-order-gateway/internal/delivery/server"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	err := server.NewServer().Run()
	if err != nil {
		panic(err)
	}

}
