package main

import (
	router "chat/routes"
)

func main() {
	r := router.Setup()
	r.Run(":3000")
}
