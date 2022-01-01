// +build ignore

package config

//go:generate strobfus -filename $GOFILE

// variables that should be obfuscated in binary

// Sleeptime is the time slept between beacons in seconds
var Sleeptime = 5

// Jitter is random extra delay to be added to the sleeptime
var Jitter = 5

// Retries is the amount of tries to keep trying in case C2 is unreachable
var Retries = 3

// UserAgent is user agent that