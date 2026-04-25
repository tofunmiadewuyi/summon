package config

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func WatchConfig(cfgCh chan<- *Config) error {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	path, err := ConfigPath()
	if err != nil {
		return err
	}

	watcher.Add(filepath.Dir(path))

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if (event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) && event.Name == path {
				// reload config
				log.Println("reloading config")
				cfg, err := Load()
				if err != nil {
					log.Println("failed to reload config:", err)
					continue
				}
				cfgCh <- cfg
			}
		}
	}

}
