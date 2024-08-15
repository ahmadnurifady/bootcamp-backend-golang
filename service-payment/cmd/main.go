package main

import (
	"fmt"
	"service-payment/internal/delivery/server"
)

func main() {

	err := server.NewServer().Run()
	if err != nil {
		fmt.Println(err)
	}
}
