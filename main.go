package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args[1:]
	w, err := NewWatcher("./")
	if err != nil {
		log.Fatal(err)
	}
	w.Run()
	defer w.Close()

	for {
		select {
		case file := <-w.files:
			if isGoFile(file) {
				pkg := pkgDir(file)
				testArgs := append([]string{pkg}, args...)
				gotest(testArgs...)
			}
		case folder := <-w.folders:
			fmt.Println("Watching path", folder)
		}
	}
}

func isGoFile(file string) bool {
	return filepath.Ext(file) == ".go"
}

func pkgDir(file string) string {
	return "./" + filepath.Dir(file)
}
