package model

import (
	"go-implant/common/communication"

	"github.com/patrickmn/go-cache"
)

var db *cache.Cache

// InitDB inits the database access
func InitDB() {

	// Create a cache with no expiration at all, and which
	// never automatically purg