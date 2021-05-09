package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"shopping/config"
	"shopping/env"
)

func Test_init(t *testing.T) {
	logEnv := "production"
	host := "http://localhost"
	port := "8080"

	os.Setenv(env.LogEnv, logEnv)
	os.Setenv(env.DBHost, host)
	os.Setenv(env.DBPort, port)
	vars := env.Init()

	cf := config.NewConfig()
	cf.Init(vars)

	assert.Equal(t, logEnv, cf.Logger.Env)
	assert.Equal(t, host, cf.Database.Host)
	assert.Equal(t, port, cf.Database.Port)
}

func Test_view(t *testing.T) {
	json := `{
  "logger": {
    "env": "development",
    "filename": "",
    "level": "debug"
  },
  "database": {
    "dialect": "",
    "host": "",
    "port": "",
    "name": "",
    "user": "",
    "password": ""
  }
}`

	cf := config.NewConfig()
	info, err := cf.View()

	assert.Equal(t, nil, err)
	assert.Equal(t, json, info)
}
