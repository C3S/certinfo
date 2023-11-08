package commands

import (
	"crypto/tls"
	"time"

	. "github.com/C3S/certinfo/internal/globals"
)

func BestBeforeHosts(confTLS *tls.Config) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	keys := sortKeys(AllHosts)
	hostQueue := make(chan map[string]*CertValidity)

	go func(q chan map[string]*CertValidity) {
		q <- daysValid(keys, AllHosts, IPversions, Timeout, confTLS, now)
	}(hostQueue)

	for r := range hostQueue {
		bestBeforeCheck(r)
	}

	close(hostQueue)
}
