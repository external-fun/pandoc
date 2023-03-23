package main

import (
	"github.com/external-fun/pandoc/api"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := api.NewDatabaseService()
	if err != nil {
		panic("Couldn't connect to database ")
	}

	statService := api.NewStatService(db)
	statService.Serve(":9090")
}
