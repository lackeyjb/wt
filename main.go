package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/logrusorgru/aurora"
)

func main() {
	fmt.Println(aurora.Cyan("wt is watching your files"))
	fmt.Printf("Type %s for help\n\n", aurora.Magenta("help"))

	testRunner(os.Args[1:])
}

func testRunner(args []string) {
	cmds := commands()
	w, err := NewWatcher("./")
	if err != nil {
		log.Fatal(err)
	}
	w.Run()
	defer w.Close()

out:
	for {
		select {
		case file := <-w.files:
			if isGoFile(file) {
				pkg := pkgDir(file)
				testArgs := append([]string{pkg}, args...)
				gotest(testArgs...)
			}
		case folder := <-w.folders:
			printWatching(folder)
		case cmd := <-cmds:
			switch cmd {
			case exit:
				break out
			case runAll:
				testArgs := append([]string{"./..."}, args...)
				gotest(testArgs...)
			case help:
				displayHelp()
			}
		}
	}
}

func isGoFile(file string) bool {
	return filepath.Ext(file) == ".go"
}

func pkgDir(file string) string {
	return "./" + filepath.Dir(file)
}

func printWatching(folder string) {
	fmt.Println("Watching path", folder)
}

func displayHelp() {
	fmt.Println(aurora.Magenta("\nInteractions:"))
	fmt.Println("  Press", aurora.White("enter").Bold(), "to run all tests")
	fmt.Println("  Press", aurora.White("q").Bold(), "to exit")
}
