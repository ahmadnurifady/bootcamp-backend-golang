package main

import (
	"fmt"
	"service-product/internal/delivery/server"
)

func main() {
	err := server.NewServer().Run()
	if err != nil {
		fmt.Println(err)
	}
}
