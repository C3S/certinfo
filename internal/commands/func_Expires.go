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

func Expires(args []string, protocol int, confTLS *tls.Config) {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: dial.DialTimeout(Timeout)},
		"tcp"+strconv.Itoa(protocol),
		args[0]+":"+strconv.Itoa(Port),
		confTLS,
	)
	if err != nil {
		if ShowErrors {
			log.Printf("%s: %s", Blue(args[0]+":"+strconv.Itoa(Port)), Orange(err))
		}
		return
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	fmt.Printf("expires: %s (IPv%d)\n", Magenta(certs[0].NotAfter.Format("02.01.2006")), protocol)
}
