package main

import (
	"ipmanlk/cnapi/api"
	"ipmanlk/cnapi/scraper"
	"ipmanlk/cnapi/sqldb"
	"sync"
)

func main() {
	sqldb.InitDB()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		scraper.Start()
	}()

	go func() {
		defer wg.Done()
		api.Start()
	}()

	wg.Wait()
}
