package main

import (
	"fmt"
	"tiny-go/api"
	"tiny-go/db"
	"tiny-go/util"
)

func main() {
	fmt.Println("Starting tiny...")

	fmt.Println("Loading configuration...")
	config := util.LoadConfiguration("config.json")

	fmt.Println("Connecting to database...")
	db.ConnectToDatabase(&config)

	fmt.Println("Starting service...")
	api.SetupService(&config)
}
