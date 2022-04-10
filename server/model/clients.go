package model

import (
	"go-implant/common/communication"

	"github.com/patrickmn/go-cache"
)

var db *cache.Cache

// InitDB inits the database access
func InitDB() {

	// Create a cache with no expiration at all, and which
	// never automatically purges expired items
	db = cache.New(cache.NoExpiration, cache.NoExpiration)
}

// Exists checks if UID exists
func Exists(UID string) bool {
	_, found := db.Get(UID)
	return found
}

// Fetch fetches the client from database. retur