package intents

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var WatchPendingTimeout = 10 * time.Minute
var WatchVerifiedTimeout = 24 * time.Hour

var WatchPendingCache *cache.Cache = nil
var WatchVerifiedCache *cache.Cache = nil

func InitIntents() {
	WatchPendingCache = cache.New(WatchPendingTimeout, 1*time.Hour)
	WatchVerifiedCache = cache.New(WatchVerifiedTimeout, 1*time.Hour)
}
