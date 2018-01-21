package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type Request struct {
	Cmd   []string `json;"cmd"`
	Chdir string   `json:"chdir,omitempty"`
}

type Handler struct {
	stdout io.Writer
	stderr io.Writer
}

func NewHandler(stdout io.Writer, stderr io.Writer) *Handler {
	return &Handler{
		stdout: stdout,
		stderr: stderr,
	}
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(421)
		fmt.Fprintln(res, err)
		return
	}

	r := &Request{}
	if err = json.Unmarshal(b, r); err != nil {
		res.WriteHeader(421)
		fmt.Fprintln(res, err)
		return
	}

	cmd := exec.Command(r.Cmd[0], r.Cmd[1:]...)
	cmd.Stdout = h.stdout
	cmd.Stderr = h.stderr
	if r.Chdir != "" {
		cmd.Dir = r.Chdir
	}
	if err = cmd.Start(); err != nil {
		res.WriteHeader(500)
		fmt.Fprintln(res, err)
		return
	}
}
