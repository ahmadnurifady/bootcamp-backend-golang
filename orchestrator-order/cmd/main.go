package main

import (
	"fmt"
	"orchestrator-order/internal/delivery/server"
)

func main() {

	err := server.NewServer().Run()
	if err != nil {
		fmt.Println(err)
	}

}
