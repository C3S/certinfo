package commands

import (
	"crypto/tls"
	"fmt"

	. "github.com/C3S/certinfo/internal/globals"
)

func ExpireHosts(confTLS *tls.Config) {
	keys := sortKeys(AllHosts)
	for _, i := range keys {
		for _, j := range IPversions {
			certs := getCerts(AllHosts[i].URL, j, AllHosts[i].Port, Timeout, confTLS)
			if certs != nil {
				fmt.Printf(
					"%-35s expires: %s (IPv%d)\n",
					Blue(AllHosts[i].URL),
					Magenta(certs[0].NotAfter.Format("02.01.2006")),
					j,
				)
			} else {
			}
		}
	}
}
