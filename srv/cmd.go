package srv

import (
	"net"

	cli "github.com/jawher/mow.cli"
)

// Version command
func Version(cmd *cli.Cmd, c net.Conn) {

	cmd.Action = func() {
		c.Write([]byte("VERSION\t v0.0.1\n"))
	}
}
