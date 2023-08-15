package commands

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/C3S/certinfo/internal/dial"
	. "github.com/C3S/certinfo/internal/globals"
)

func BestBefore(args []string, protocol int, confTLS *tls.Config) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: dial.DialTimeout(Timeout)},
		"tcp"+strconv.Itoa(protocol),
		args[0]+":"+strconv.Itoa(Port),
		confTLS,
	)
	if err != nil {
		if ShowErrors {
			log.Printf(
				"%s: %s %s",
				Blue(args[0]+":"+strconv.Itoa(Port)),
				Red("Error during dial:"),
				Orange(err),
			)
		}
		return
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
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
}
