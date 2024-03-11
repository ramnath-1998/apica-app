package main

import (
	"sync"

	"github.com/ramnath.1998/apica-app/handlers"
	"github.com/ramnath.1998/apica-app/routes"
)

func init() {

}

func main() {
	cache := handlers.CacheStruct{}
	var wg sync.WaitGroup
	wg.Add(1)
	cache.InitializeCache(1024)
	go cache.RemoveExpired()
	routes.RunRoutes()
	wg.Wait()

}
