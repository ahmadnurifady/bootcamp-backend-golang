package main

import "service-user/internal/delivery/server"

func main() {
	err := server.NewServer().Run()
	if err != nil {
		return
	}

	//result, err := outbond.GetUserValidation("usr-0021")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(*result)
}
