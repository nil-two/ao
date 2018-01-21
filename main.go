package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	cliName    = "ao"
	cliVersion = "0.1.0"
)

type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	serve  bool
}

func NewCLI(stdin io.Reader, stdout io.Writer, stderr io.Writer) *CLI {
	return &CLI{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (c *CLI) printUsage() {
	fmt.Fprintf(c.stderr, `
usage: %[1]s <command> [...]
await order and execute it.

commands:
  %[1]s a|await [-port=<port>]                  # await order
  %[1]s o|order [-port=<port>] <cmd> [<arg(s)>] # order to execute command
  %[1]s h|help                                  # display usage
  %[1]s v|version                               # display version
`[1:], cliName)
}

func (c *CLI) printVersion() {
	fmt.Fprintf(c.stdout, "%s\n", cliVersion)
}

func (c *CLI) printError(err interface{}) {
	fmt.Fprintf(c.stderr, "%s: %s\n", cliName, err)
}

func (c *CLI) await(args []string) int {
	f := flag.NewFlagSet("await", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	port := f.Int("port", 60080, "port number")
	if err := f.Parse(args); err != nil {
		c.printError(err)
		return 2
	}

	addr := fmt.Sprintf(":%d", *port)
	handler := NewHandler(c.stdout, c.stderr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		c.printError(err)
		return 1
	}
	return 0
}

func (c *CLI) order(args []string) int {
	f := flag.NewFlagSet("await", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	port := f.Int("port", 60080, "port number")
	if err := f.Parse(args); err != nil {
		c.printError(err)
		return 2
	}
	cmd := f.Args()

	b, err := json.Marshal(&Request{Cmd: cmd})
	if err != nil {
		c.printError(err)
		return 1
	}

	url := fmt.Sprintf("http://localhost:%d/", *port)
	res, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		c.printError(err)
		return 1
	}
	defer res.Body.Close()

	if _, err = io.Copy(c.stdout, res.Body); err != nil {
		c.printError(err)
		return 1
	}
	return 0
}

func (c *CLI) help(args []string) int {
	c.printUsage()
	return 0
}

func (c *CLI) version(args []string) int {
	c.printVersion()
	return 0
}

func (c *CLI) Run(args []string) int {
	if len(args) < 1 {
		c.printUsage()
		return 2
	}

	cmd := args[0]
	switch cmd {
	case "a", "await":
		return c.await(args[1:])
	case "o", "order":
		return c.order(args[1:])
	case "h", "help":
		return c.help(args[1:])
	case "v", "version":
		return c.version(args[1:])
	}

	c.printError(cmd + ": no such command")
	return 1
}

func main() {
	c := NewCLI(os.Stdin, os.Stdout, os.Stderr)
	e := c.Run(os.Args[1:])
	os.Exit(e)
}
