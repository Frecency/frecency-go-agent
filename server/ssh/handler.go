// +build !windows

// is this needed???

package ssh

import (
	"go-implant/client/config"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/fatih/