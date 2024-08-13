// package internal

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/fsnotify/fsnotify"
// )

// // WatchConfig watches the INI file for changes and invokes the provided callback with the new configuration.
// func WatchConfig(ctx context.Context, filename string, onChange func(*Config)) error {
// 	if _, err := os.Stat(filename); os.IsNotExist(err) {
// 		return fmt.Errorf("file does not exist: %v", err)
// 	}

// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		return fmt.Errorf("error creating watcher: %v", err)
// 	}
// 	defer watcher.Close()

// 	if err := watcher.Add(filename); err != nil {
// 		return fmt.Errorf("error adding file to watcher: %v", err)
// 	}

// 	for {
// 		select {
// 		case event := <-watcher.Events:
// 			if event.Op&fsnotify.Write == fsnotify.Write {
// 				cfg, err := LoadConfig(filename)
// 				if err != nil {
// 					log.Printf("Error reloading config: %v", err)
// 					continue
// 				}
// 				onChange(cfg)
// 			}
// 		case err := <-watcher.Errors:
// 			log.Printf("Watcher error: %v", err)
// 			return err
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		}
// 	}
// }

package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

// WatchConfig watches the INI file for changes and invokes the provided callback with the new configuration.
func WatchConfig(ctx context.Context, filename string, onChange func(*Config)) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %v", err)
	}
	defer func() {
		if closeErr := watcher.Close(); closeErr != nil {
			logrus.Errorf("Error closing watcher: %v", closeErr)
		}
	}()

	if err := watcher.Add(filename); err != nil {
		return fmt.Errorf("error adding file to watcher: %v", err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				cfg, err := LoadConfig(filename)
				if err != nil {
					logrus.Errorf("Error reloading config: %v", err)
					continue
				}
				onChange(cfg)
			}
		case err := <-watcher.Errors:
			logrus.Errorf("Watcher error: %v", err)
			return err
		case <-ctx.Done():
			logrus.Infof("Context done, stopping watcher.")
			return ctx.Err()
		}
	}
}

