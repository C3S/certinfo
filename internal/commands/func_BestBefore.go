package commands

import (
	"crypto/tls"
	"time"

	. "github.com/C3S/certinfo/internal/globals"
)

func BestBefore(args []string, protocol int, confTLS *tls.Config) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	allHosts := make(map[string]*Host)
	thisHost := &Host{
		URL:  args[0],
		Port: Port,
	}
	allHosts[args[0]] = thisHost

	thisKey := []string{args[0]}

	thisProtocol := [2]int{protocol, 0}

	bestBeforeCheck(
		thisKey,
		allHosts,
		thisProtocol,
		Timeout,
		confTLS,
		now,
	)
}
