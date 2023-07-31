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
	Values []lokiStreamLine  `json:"values"`
}

type lokiStreamLine timedLine

func (tl *lokiStreamLine) MarshalJSON() ([]byte, error) {
	bs, err := json.Marshal(timedLine(*tl))
	if err != nil {
		return nil, err
	}
	return json.Marshal([]string{
		fmt.Sprintf("%d", tl.Time.UnixNano()),
		string(bs),
	})
}

func postLines(labels map[string]string, lines []timedLine, cli *CLI) error {
	lokiLines := make([]lokiStreamLine, len(lines))
	for i, tl := range lines {
		lokiLines[i] = lokiStreamLine(tl)
	}
	re := request{
		Streams: []stream{
			{
				Labels: labels,
				Values: lokiLines,
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
