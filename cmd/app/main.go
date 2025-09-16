package main

import (
	"log"
	"sif/src/apps"
	"sif/src/config"
	"time"

	"github.com/socious-io/goaccount"
	"github.com/socious-io/gomq"
	"github.com/socious-io/gopay"
	database "github.com/socious-io/pkg_database"
)

func main() {
	config.Init("config.yml")
	database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 50,
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

	//Initializing GoMQ Library
	gomq.Setup(gomq.Config{
		Url:        config.Config.Nats.Url,
		Token:      config.Config.Nats.Token,
		ChannelDir: "fund",
	})
	gomq.Connect()

	if err := goaccount.Setup(config.Config.GoAccounts); err != nil {
		log.Fatalf("goaccount error %v", err)
	}

	apps.Serve()
}
