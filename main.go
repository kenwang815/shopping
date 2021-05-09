package main

import (
	"time"

	preparation "shopping/init"
	"shopping/repository"
	"shopping/rest"
	"shopping/service"
	"shopping/utils/log"
)

func main() {
	// Init config
	cf := preparation.Init()
	cf.View()

	// Init repository
	e, err := repository.NewEngine(cf)
	if err != nil {
		log.Error(err)
	}

	e.Database.SetPool(10, 100, time.Hour)

	// Init service
	err = service.Init(cf, e)
	if err != nil {
		log.Error(err)
	}

	// Init rest
	router := rest.Init()
	router.Run(":8080")
}
