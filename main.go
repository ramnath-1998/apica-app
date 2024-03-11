package main

import (
	"github.com/ramnath.1998/apica-app/handlers"
	"github.com/ramnath.1998/apica-app/routes"
)

func init() {

}

func main() {

	handlers.InitializeCache(4)
	go handlers.RemoveExpired()
	routes.RunRoutes()

}
