package main

import (
	router "github.com/mikejeffers/chat-api-go/routes"
)

func main() {
	r := router.Setup()
	r.Run(":3000")
}
