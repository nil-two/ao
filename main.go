package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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

	server := NewServer(*port, c.stdout, c.stderr)
	if err := server.Serve(); err != nil {
		c.printError(err)
		return 1
	}
	return 0
}

func (c *CLI) order(args []string) int {
	f := flag.NewFlagSet("order", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	port := f.Int("port", 60080, "port number")
	if err := f.Parse(args); err != nil {
		c.printError(err)
		return 2
	}
	if f.NArg() < 1 {
		c.printError("no input cmd")
		return 2
	}

	cmd := f.Args()

	client := NewClient(*port, c.stdout)
	if err := client.Order(cmd); err != nil {
		c.printError(err)
		return 1
	}
	return 0
}

func (c *CLI) help(args []string) int {
	f := flag.NewFlagSet("help", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	if err := f.Parse(args); err != nil {
		c.printError(err)
		return 2
	}

	c.printUsage()
	return 0
}

func (c *CLI) version(args []string) int {
	f := flag.NewFlagSet("version", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	if err := f.Parse(args); err != nil {
		c.printError(err)
		return 2
	}

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
	default:
		c.printError(cmd + ": no such command")
		return 1
	}
}

func main() {
	c := NewCLI(os.Stdin, os.Stdout, os.Stderr)
	e := c.Run(os.Args[1:])
	os.Exit(e)
}
