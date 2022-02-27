package communication

import (
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// BeaconResponse struct which contains
// an array of users
type BeaconResponse struct {
	Commands []Command `json:"Commands"`
}

// Command struct which contains
// an command and arguments
type Command struct {
	Command string   `json