package commands

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	. "github.com/C3S/certinfo/internal/globals"
)

func bestBeforeCheck(
	keys []string,
	allHosts map[string]*Host,
	IPversions [2]int,
	timeout int,
	confTLS *tls.Config,
	now time.Time,
) {
	for _, i := range keys {
		for _, j := range IPversions {
			if j != 4 && j != 6 {
				continue
			} else {
				certs := getCerts(allHosts[i].URL, j, allHosts[i].Port, timeout, confTLS)
				if certs != nil {
					certExpires := certs[0].NotAfter
					daysValid := int(certExpires.Sub(now).Hours() / 24)
					if daysValid > Days {
						fmt.Printf(
							"%-35s (IPv%d): expires %-44s %s",
							Blue(allHosts[i].URL),
							j,
							Green(certs[0].NotAfter.Format("02.01.2006"))+", in "+Green(strconv.Itoa(daysValid))+" days",
							Green("-- ok!\n"),
						)
						continue
					} else if daysValid < 0 {
						fmt.Printf(
							"%-35s (IPv%d): expired %-44s %s",
							Blue(allHosts[i].URL),
							j,
							Red(certs[0].NotAfter.Format("02.01.2006"))+", "+Red(strconv.Itoa(daysValid))+" days ago",
							Red("-- red alert!\n"),
						)
						continue
					} else {
						fmt.Printf(
							"%-35s (IPv%d): expires %-44s %s",
							Blue(allHosts[i].URL),
							j,
							Orange(certs[0].NotAfter.Format("02.01.2006"))+", in "+Orange(strconv.Itoa(daysValid))+" days",
							Orange("-- please renew!\n"),
						)
						continue
					}
				} else {
				}
			}
		}
	}
}
