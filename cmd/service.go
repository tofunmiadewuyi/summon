package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/tofunmiadewuyi/summon/internal/accessibility"
	"github.com/tofunmiadewuyi/summon/internal/config"
	"github.com/tofunmiadewuyi/summon/internal/hotkey"
	"github.com/tofunmiadewuyi/summon/internal/service"
)

func run() {
	runtime.LockOSThread()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// register combos from config
	hotkey.Register(cfg)

	// ensure they have accesibility allowed for summon, program does not work without it.
	// instead of blocking and restarting (os.Exit), we show the prompt and retry the tap
	// in the background — this works correctly from both terminal and launchd contexts.
	if !hotkey.Start() {
		accessibility.RequestPermission()
		log.Println("accessibility permission required — grant access in System Settings")
		go func() {
			for !hotkey.IsRunning() {
				time.Sleep(3 * time.Second)
				hotkey.Start()
			}
			log.Println("event tap active ✨")
		}()
	}

	// watch the config for hot reloads
	cfgCh := make(chan *config.Config)
	go config.WatchConfig(cfgCh)

	go func() {
		for cfg := range cfgCh {
			hotkey.Register(cfg)
			hotkey.Start()
		}
	}()

	fmt.Println("ready to summon ✨")
	hotkey.RunMainLoop()
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
