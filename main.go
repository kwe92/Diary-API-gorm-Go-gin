package main

import (
	"fmt"
	"journal_api/router"
	"journal_api/utility"
)

// TODO: Review function / method comments

const address = ":8080"

func main() {
	utility.LoadEnv()

	utility.LoadDatabase()

	router := router.SetupRouter()

	router.Run(address)

	fmt.Println("Server running on port 8000")

}
