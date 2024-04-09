package main

import (
	"ipmanlk/cnapi/api"
	"ipmanlk/cnapi/nosqldb"
	"ipmanlk/cnapi/scraper"
	"sync"
)

func main() {
	err := nosqldb.InitDB()
	if err != nil {
		panic(err)
	}

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
