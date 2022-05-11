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
	fmt.Printf("CurrentUser: %s\n",