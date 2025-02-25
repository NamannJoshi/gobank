package main

import (
	"fmt"
	"gobank/api"
	"log"
)

func main() {
	fmt.Println("Jai baba ri")
	fmt.Println("hey are you suar, because i want to get dirty with you")

	store, err := api.NewPostgreStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal("Init error:", err)
	}

	fmt.Printf("%v+\n", store)

	server := api.NewServerApi(":3000", store)
	server.Run()
}