package init

import (
	"shopping/config"
	"shopping/env"
	"shopping/utils/log"
)

func Init() *config.Config {
	// Init env
	vars := env.Init()

	// Init config
	cf := config.NewConfig()
	cf.Init(vars)

	// Init logger
	log.Init(cf.Logger.Env, cf.Logger.Filename, cf.Logger.Level)

	return cf
}
