package commands

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/C3S/certinfo/internal/dial"
	. "github.com/C3S/certinfo/internal/globals"
)

func ExpireHosts(confTLS *tls.Config) {
	keys := sortKeys(AllHosts)
	for _, i := range keys {
		for _, j := range IPversions {
			conn, err := tls.DialWithDialer(
				&net.Dialer{Timeout: dial.DialTimeout(Timeout)},
				"tcp"+strconv.Itoa(j),
				AllHosts[i].URL+":"+strconv.Itoa(AllHosts[i].Port),
				confTLS,
			)
			if err != nil {
				if ShowErrors {
					log.Printf(
						"%-35s %s",
						Blue(AllHosts[i].URL+":"),
						Orange(err),
					)
				}
				continue
			}
			defer conn.Close()
			certs := conn.ConnectionState().PeerCertificates
			fmt.Printf(
				"%-35s expires: %s (IPv%d)\n",
				Blue(AllHosts[i].URL),
				Magenta(certs[0].NotAfter.Format("02.01.2006")),
				j,
			)
		}
	}
}
