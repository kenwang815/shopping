package main

import (
	"time"

	preparation "shopping/init"
	"shopping/repository"
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
}
