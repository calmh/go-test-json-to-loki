package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type request struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Labels map[string]string `json:"stream"`
	Values []timedLine       `json:"values"`
}

func postLines(labels map[string]string, lines []timedLine, cli *CLI) error {
	re := request{
		Streams: []stream{
			{
				Labels: labels,
				Values: lines,
			},
		},
	}
	bs, err := json.Marshal(re)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, cli.LokiURL, bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if cli.LokiUser != "" {
		req.SetBasicAuth(cli.LokiUser, cli.LokiPassword)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		io.Copy(os.Stderr, resp.Body)
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
