package main

import (
	"fmt"
	"gobank/api"
)

func main() {
	fmt.Println("Jai baba ri")
	fmt.Println("hey are you suar, because i want to get dirty with you")

	server := api.NewServerApi(":3000")
	server.Run()
}