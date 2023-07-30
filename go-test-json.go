package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type timedLine struct {
	Time   time.Time
	Output string
	raw    string
}

func (tl *timedLine) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string{
		fmt.Sprintf("%d", tl.Time.UnixNano()),
		tl.raw,
	})
}

func readStdinInto(ch chan<- timedLine) {
	br := bufio.NewReader(os.Stdin)
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			break
		}

		var tl timedLine
		tl.raw = line
		if err := json.Unmarshal([]byte(line), &tl); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling JSON: %v\n", err)
			continue
		}
		ch <- tl
	}
}
