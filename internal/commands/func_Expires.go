package commands

import (
	"crypto/tls"
	"fmt"

	. "github.com/C3S/certinfo/internal/globals"
)

func Expires(args []string, protocol int, confTLS *tls.Config) {
	certs := getCerts(args[0], protocol, Port, Timeout, confTLS)
	if certs != nil {
		fmt.Printf("expires: %s (IPv%d)\n", Magenta(certs[0].NotAfter.Format("02.01.2006")), protocol)
	} else {}
}
