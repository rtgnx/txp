package cmd

import (
	"flag"
	"fmt"
	"net"

	cli "github.com/jawher/mow.cli"
)

// Handler type
type Handler func(args []string, conn net.Conn) error

// MultiHandler muliple handler parser
func MultiHandler(handlers map[string]Handler) Handler {
	return func(args []string, conn net.Conn) error {
		if len(args) < 1 {
			return fmt.Errorf("expected at least one argument")
		}

		if h, ok := handlers[args[0]]; ok {
			return h(args, conn)
		} else if h, ok := handlers["default"]; ok {
			return h(args, conn)
		}

		return nil
	}
}

// CommandHandler type
type CommandHandler func(*cli.Cmd, net.Conn)

// Command handler
func Command(name, description string, cmd CommandHandler) Handler {
	return func(args []string, conn net.Conn) error {
		app := cli.App(name, description)
		app.ErrorHandling = flag.ContinueOnError
		cmd(app.Cmd, conn)
		return app.Run(args)
	}
}
