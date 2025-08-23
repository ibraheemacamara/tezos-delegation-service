package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/ibraheemacara/tezos-delegation-service/api"
	"github.com/ibraheemacara/tezos-delegation-service/config"
	"github.com/ibraheemacara/tezos-delegation-service/db"
	delegationswatcher "github.com/ibraheemacara/tezos-delegation-service/delegations_watcher"
	"github.com/ibraheemacara/tezos-delegation-service/httpclient"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the config file, if empty string defaults will be used")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	db, err := db.InitDB(cfg)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	httpClient := httpclient.NewHttpClient(4 * time.Second)
	delegationsWatcher := delegationswatcher.NewDelegationsWatcher(cfg, httpClient, db)
	go delegationsWatcher.Start()

	api.StartServer(cfg, db)
}
