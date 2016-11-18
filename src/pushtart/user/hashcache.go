package user

import (
	"crypto/sha1"

	lru "github.com/hashicorp/golang-lru"
)

// To speed up web serving (where the hash needs to be verified for every request)
// we can avoid the overhead of the bcrypt check by caching good username-SHA(password) combinations.
// we could cache the password in memory, but I thought instead doing the SHA is just one more level
// an attacker with root / memory-read has to go through.

var lookupCache *lru.ARCCache

const cacheExpiryTime = 60 * 180 //180 minutes
const hashCacheSize = 20
const salt = "ShittySalt!543*"

func init() {
	var err error
	lookupCache, err = lru.NewARC(hashCacheSize)
	if err != nil {
		panic(err)
	}
}

func isValidFromCacheHit(username, pass string) bool {
	cacheVal, cacheHit := lookupCache.Get(username)
	if !cacheHit {
		return false
	}
	return cacheVal.([sha1.Size]byte) == sha1.Sum([]byte(pass+salt+username))
}

func cacheCorrectAuthEntry(username, pass string) {
	lookupCache.Add(username, sha1.Sum([]byte(pass+salt+username)))
}
