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

// On wind