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
func BeaconHandler(w http.ResponseWriter, r *http.Request