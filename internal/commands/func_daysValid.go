package commands

import (
	"crypto/tls"
	"time"

	. "github.com/C3S/certinfo/internal/globals"
)

func daysValid(
	keys []string,
	allHosts map[string]*Host,
	IPversions [2]int,
	timeout int,
	confTLS *tls.Config,
	now time.Time,
) map[string]*CertValidity {
	daysValidLeft := make(map[string]*CertValidity)
	for _, i := range keys {
		for _, j := range IPversions {
			if j != 4 && j != 6 {
				continue
			} else {
				certs := getCerts(allHosts[i].URL, j, allHosts[i].Port, timeout, confTLS)
				if certs != nil {
					certExpires := certs[0].NotAfter
					daysValidLeft[i] = &CertValidity{
					  URL:	allHosts[i].URL,
					  Port: allHosts[i].Port,
					  Protocol: j,
					  Certificate: certs,
					  DaysLeft: int(certExpires.Sub(now).Hours() / 24),
					}
				} else {
				}
			}
		}
	}
	return daysValidLeft
}
