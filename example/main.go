package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/rtgnx/txp/client"
	"github.com/rtgnx/txp/cmd"
	"github.com/rtgnx/txp/srv"

	cli "github.com/jawher/mow.cli"
)

var commands = map[string]cmd.Handler{
	"VER":        cmd.Command("VER", "version", version),
	"ECHO":       cmd.Command("ECHO", "echo", echo),
	"MSG":        cmd.Command("MSG", "message", message),
	"disconnect": cmd.Command("disconnect", "on disconnect", disconnect),
	"connect":    cmd.Command("connect", "on connect", connect),
}

var clients = map[string]net.Conn{}

func message(cmd *cli.Cmd, c net.Conn) {
	msg := cmd.StringArg("MSG", "", "")
	cmd.Action = func() {
		log.Printf("%s", *msg)
	}
}

func disconnect(cmd *cli.Cmd, c net.Conn) {
	cmd.Action = func() {
		delete(clients, c.RemoteAddr().String())
	}
}

func connect(cmd *cli.Cmd, c net.Conn) {
	cmd.Action = func() {
		c.Write([]byte("Welcome to 1337 Server\n"))
		clients[c.RemoteAddr().String()] = c
	}
}

func version(cmd *cli.Cmd, c net.Conn) {

	cmd.Action = func() {
		c.Write([]byte("VERSION\t v0.0.1\n"))
	}
}

func echo(cmd *cli.Cmd, c net.Conn) {
	var (
		msg = cmd.StringArg("MSG", "", "message")
	)

	cmd.Action = func() {
		for _, c := range clients {
			c.Write([]byte("MSG " + *msg))
		}
	}
}

func main() {
	app := cli.App("example", "example client/server")

	app.Command("serve", "", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			server := srv.Server{Commands: commands}
			server.Start(":1337")
		}
	})

	app.Command("client", "", func(cmd *cli.Cmd) {
		addr := cmd.StringArg("ADDR", "", "address")
		cmd.Action = func() {

			c, _ := client.New(*addr, commands)
			c.Connect(context.Background())
		}
	})

	app.Run(os.Args)

}
