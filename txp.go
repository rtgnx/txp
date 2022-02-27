package txp

import (
	"github.com/rtgnx/txp/client"
	"github.com/rtgnx/txp/cmd"
	"github.com/rtgnx/txp/srv"
)

// NewServer instance
func NewServer(commands map[string]cmd.Handler) *srv.Server {
	return &srv.Server{Commands: commands}
}

// NewClient instance
func NewClient(addr string, cmds map[string]cmd.Handler) (*client.Client, error) {
	return client.New(addr, cmds)
}
