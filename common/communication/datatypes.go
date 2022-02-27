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
