package commands

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"time"

	. "github.com/C3S/certinfo/internal/globals"
)

func BestBefore(args []string, protocol int, confTLS *tls.Config) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	certs := getCerts(args[0], protocol, Port, Timeout, confTLS)
	if certs != nil {
		certExpires := certs[0].NotAfter
		daysValid := int(certExpires.Sub(now).Hours() / 24)
		if daysValid > Days {
			fmt.Printf(
				"will expire in %s days (IPv%d) %s",
				Green(strconv.Itoa(daysValid)),
				protocol,
				Green(" -- ok!\n"),
			)
			return
		} else if daysValid < 0 {
			fmt.Printf(
				"expired %s days ago (IPv%d) %s",
				Red(strconv.Itoa(daysValid)),
				protocol,
				Red(" -- red alert!\n"),
			)
			os.Exit(1)
		} else {
			fmt.Printf(
				"expires in %s days (IPv%d) %s",
				Orange(strconv.Itoa(daysValid)),
				protocol,
				Orange(" -- please renew!\n"),
			)
			os.Exit(1)
		}
	} else {
	}
}
