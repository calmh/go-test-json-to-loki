package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type timedLine struct {
	// fields from go test -json
	Time    time.Time
	Output  string  `json:",omitempty"`
	Action  string  `json:",omitempty"`
	Test    string  `json:",omitempty"`
	Package string  `json:",omitempty"`
	Elapsed float64 `json:",omitempty"`

	// extra fields
	Run string
}

func readStdinInto(ch chan<- timedLine) {
	br := bufio.NewReader(os.Stdin)
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			break
		}

		var tl timedLine
		if err := json.Unmarshal([]byte(line), &tl); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling JSON: %v\n", err)
			continue
		}
		ch <- tl
	}
}
