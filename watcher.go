package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	*fsnotify.Watcher
	files   chan string
	folders chan string
}

func NewWatcher(path string) (*watcher, error) {
	folders, err := subfolders(path)
	if err != nil {
		return nil, err
	}
	if len(folders) == 0 {
		return nil, errors.New("No folders to watch.")
	}

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &watcher{Watcher: fsWatcher}
	w.files = make(chan string, 10)
	w.folders = make(chan string, len(folders))

	for _, folder := range folders {
		w.addFolder(folder)
	}

	return w, nil
}

func (w *watcher) addFolder(folder string) {
	err := w.Add(folder)
	if err != nil {
		log.Println("Error watching: ", folder, err)
	}
	w.folders <- folder
}

func (w *watcher) Run() {
	go func() {
		for {
			select {
			case event := <-w.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					file, _ := os.Stat(event.Name)
					if file.IsDir() && !shouldIgnoreDir(filepath.Base(event.Name)) {
						w.addFolder(event.Name)
					} else {
						w.files <- event.Name
					}
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					w.files <- event.Name
				}
			case err := <-w.Errors:
				log.Println(err)
			}
		}
	}()
}

func subfolders(startPath string) ([]string, error) {
	var paths []string
	err := filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if shouldIgnoreDir(name) {
				return filepath.SkipDir
			}

			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func shouldIgnoreDir(name string) bool {
	return (strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")) && name != "." && name != ".."
}
