package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	port   int
	stdout io.Writer
}

func NewClient(port int, stdout io.Writer) *Client {
	return &Client{
		port:   port,
		stdout: stdout,
	}
}

func (c *Client) Order(cmd []string) error {
	b, err := json.Marshal(&Request{Cmd: cmd})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:%d/", c.port)
	res, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if _, err = io.Copy(c.stdout, res.Body); err != nil {
		return err
	}
	return nil
}
