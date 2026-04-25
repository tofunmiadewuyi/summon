package main

import (
	"log"

	"github.com/tofunmiadewuyi/summon/internal/config"
	"github.com/tofunmiadewuyi/summon/internal/hotkey"
	"github.com/tofunmiadewuyi/summon/internal/service"
)

func run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	hotkey.Register(cfg)

	cfgCh := make(chan *config.Config)
	go config.WatchConfig(cfgCh)

	go func() {
		for cfg := range cfgCh {
			hotkey.Register(cfg)
		}
	}()

	log.Println("ready to summon ✨")
	//block forever
	select {}
}

func start() {
	service.LaunchdStart()
}

func status() {
	service.LaunchdStatus()
}

func stop() {
	service.LaunchdStop()
}
