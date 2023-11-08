package commands

import (
	"crypto/tls"
	"time"

	. "github.com/C3S/certinfo/internal/globals"
	// "github.com/schollz/progressbar/v3"
)

func BestBeforeHosts(confTLS *tls.Config) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	keys := sortKeys(AllHosts)
	hostQueue := make(chan map[string]*CertValidity, 5)
	// bar := progressbar.Default(int64(len(keys)))

	go func(q chan map[string]*CertValidity) {
		// for j := 0; j < len(keys); j++ {
			// thisKey := []string{keys[j]}
			// q <- daysValid(thisKey, AllHosts, IPversions, Timeout, confTLS, now)
			// bar.Add(1)
		// }
		q <- daysValid(keys, AllHosts, IPversions, Timeout, confTLS, now)
	}(hostQueue)

	for r := range hostQueue {
		bestBeforeCheck(r)
	}

	close(hostQueue)
}

