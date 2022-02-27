package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/rtgnx/txp/cmd"
)

// Client object
type Client struct {
	conn     net.Conn
	Commands map[string]cmd.Handler
}

// New client instance
func New(addr string, cmds map[string]cmd.Handler) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	return &Client{conn: conn, Commands: cmds}, err
}

// Command send command to server
func (c *Client) Command(args []string) error {
	_, err := c.conn.Write([]byte(strings.Join(append(args, "\n"), " ")))
	return err
}

// Connect client to server
func (c *Client) Connect(ctx context.Context) error {

	r := bufio.NewReader(c.conn)
	handler := cmd.MultiHandler(c.Commands)

	for {
		line, err := r.ReadBytes('\n')
		if err == io.EOF {
			return fmt.Errorf("connection terminated")
		}

		go func() {
			if err := handler(strings.Fields(string(line)), c.conn); err != nil {
				log.Println(err)
			}
		}()
	}

}
