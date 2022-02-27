package srv

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"

	"github.com/rtgnx/txp/cmd"
)

// Context passed to each command
type Context struct {
	Cmd  string
	Args []string
	conn net.Conn
}

// Write message back to the client
// failed write means connection has been terminated
// command should be written in request / response style
// or cleanup and terminate on first write failure
func (c *Context) Write(args ...string) error {
	_, err := c.conn.Write([]byte(strings.Join(args, " ") + "\n"))
	return err
}

// Server object
type Server struct {
	// if connect and disconnect commands are present
	// they will be called on iniiated connection and upon termination
	Commands map[string]cmd.Handler
}

// Start TCP server
func (s *Server) Start(addr string) error {
	fd, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	for {
		conn, err := fd.Accept()

		if err != nil {
			continue
		}

		go s.connHandler(conn)
	}
}

func (s *Server) connHandler(c net.Conn) {
	buf := bufio.NewReader(c)
	// On connect hook
	if h, ok := s.Commands["connect"]; ok {
		h([]string{c.RemoteAddr().String()}, c)
	}

	for {
		line, err := buf.ReadBytes('\n')

		if err == io.EOF {
			// on disconnect hook
			if h, ok := s.Commands["disconnect"]; ok {
				h([]string{c.RemoteAddr().String()}, c)
			}
			return
		}

		args := strings.Fields(string(line))
		if err := cmd.MultiHandler(s.Commands)(args, c); err != nil {
			log.Println(err)
		}
	}
}
