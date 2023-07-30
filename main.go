package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
)

type CLI struct {
	LokiURL      string   `help:"Loki URL" required:"true" env:"LOKI_URL"`
	LokiUser     string   `help:"Loki user" env:"LOKI_USER"`
	LokiPassword string   `help:"Loki password" env:"LOKI_PASSWORD"`
	Labels       []string `help:"Labels to add to all lines (key=value)" env:"LOKI_LABELS" short:"l"`
}

func main() {
	var cli CLI
	kong.Parse(&cli)

	labels := make(map[string]string)
	for _, l := range cli.Labels {
		k, v, ok := strings.Cut(l, "=")
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: invalid label %q\n", l)
			os.Exit(1)
		}
		labels[k] = v
	}

	in := make(chan timedLine)
	go func() {
		readStdinInto(in)
		close(in)
	}()

	postTimer := time.NewTicker(5 * time.Second)

	var lines []timedLine
	exitCode := 0
loop:
	for {
		select {
		case tl, ok := <-in:
			if !ok {
				break loop
			}
			lines = append(lines, tl)
			if tl.Output != "" {
				fmt.Print(tl.Output)
			}
			if tl.Action == "fail" {
				// A test failed; propagate the exit code.
				exitCode = 1
			}

		case <-postTimer.C:
			if len(lines) > 0 {
				if err := postLines(labels, lines, &cli); err != nil {
					fmt.Fprintf(os.Stderr, "Error: posting to Loki: %v\n", err)
				}
				lines = lines[:0]
			}
		}
	}

	if len(lines) > 0 {
		if err := postLines(labels, lines, &cli); err != nil {
			fmt.Fprintf(os.Stderr, "Error: posting to Loki: %v\n", err)
		}
	}

	os.Exit(exitCode)
}
