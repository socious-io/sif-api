package main

import (
	"sif/src/apps/workers"
	"sif/src/config"
	"time"

	"github.com/socious-io/gomq"
	database "github.com/socious-io/pkg_database"
)

func main() {

	config.Init("config.yml")
	database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})

	gomq.Setup(gomq.Config{
		Url:        config.Config.Nats.Url,
		Token:      config.Config.Nats.Token,
		ChannelDir: "fund",
		Consumers:  map[string]func(interface{}){},
	})

	workers.RegisterConsumers()

	gomq.Init()
}
