package model

import (
	"go-implant/common/communication"

	"github.com/patrickmn/go-cache"
)

var db *cache.Cache

// InitDB ini