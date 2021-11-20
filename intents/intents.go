package intents

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var WatchPendingCache *cache.Cache = nil
var WatchVerifiedCache *cache.Cache = nil

func InitIntents() {
	WatchPendingCache = cache.New(10*time.Minute, 1*time.Hour)
	WatchVerifiedCache = cache.New(24*time.Hour, 1*time.Hour)
}
