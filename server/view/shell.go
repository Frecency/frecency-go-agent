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

var lette