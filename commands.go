package main

import (
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

type command int

const (
	unknown command = iota + 1
	help
	exit
	runAll
)

func commands() <-chan command {
	cmds := make(chan command, 1)
	rl, err := readline.New("")
	if err != nil {
		log.Fatal(err)
	}
	// defer rl.Close()

	go func() {
		for {
			line, err := rl.Readline()
			if err == io.EOF {
				cmds <- exit
				break
			} else if err != nil {
				log.Fatal(err)
			}

			cmds <- normalizeCommand(line)
			if err := readline.AddHistory(line); err != nil {
				log.Fatal(err)
			}
		}
	}()

	return cmds
}

func normalizeCommand(line string) command {
	cmd := strings.ToLower(strings.TrimSpace(line))
	switch cmd {
	case "exit", "e", "x", "quit", "q":
		return exit
	case "all", "a", "":
		return runAll
	case "help", "h", "?":
		return help
	default:
		return unknown
	}
}
