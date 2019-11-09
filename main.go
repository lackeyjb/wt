package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/logrusorgru/aurora"
)

func main() {
	os.Exit(gotest(os.Args[1:]))
}

func gotest(args []string) int {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	r, w := io.Pipe()
	defer w.Close()

	args = append([]string{"test"}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = w
	cmd.Stdout = w
	cmd.Env = os.Environ()

	go consume(&wg, r)

	if err := cmd.Run(); err != nil {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			return ws.ExitStatus()
		}
		return 1
	}

	return 0
}

func consume(wg *sync.WaitGroup, r io.Reader) {
	defer wg.Done()
	reader := bufio.NewReader(r)
	for {
		line, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			log.Print(err)
			return
		}
		parse(string(line))
	}
}

func parse(line string) {
	trimmed := strings.TrimSpace(line)
	color := trimmed

	switch {
	case strings.HasPrefix(trimmed, "--- PASS"):
		fallthrough
	case strings.HasPrefix(trimmed, "ok"):
		fallthrough
	case strings.HasPrefix(trimmed, "PASS"):
		color = colorizeLine(line, aurora.GreenFg)
	case strings.HasPrefix(trimmed, "--- FAIL"):
		fallthrough
	case strings.HasPrefix(trimmed, "FAIL"):
		color = colorizeLine(line, aurora.RedFg)
	}
	fmt.Printf("%s\n", color)
}

func colorizeLine(line string, color aurora.Color) string {
	return aurora.Colorize(line, color).String()
}
