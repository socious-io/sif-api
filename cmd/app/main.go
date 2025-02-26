package main

import (
	"log"
	"socious/src/apps"
	"socious/src/config"
	"time"

	"github.com/socious-io/gopay"
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

	if err := gopay.Setup(gopay.Config{
		DB:     database.GetDB(),
		Prefix: "gopay",
		Chains: config.Config.Payment.Chains,
		Fiats:  config.Config.Payment.Fiats,
	}); err != nil {
		log.Fatalf("gopay error %v", err)
	}

	apps.Serve()
}
