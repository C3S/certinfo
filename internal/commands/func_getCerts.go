package commands

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"strconv"
	"time"

	. "github.com/C3S/certinfo/internal/globals"
)

func getCerts(
	url string,
	protocol int,
	port int,
	timeout int,
	confTLS *tls.Config,
) []*x509.Certificate {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: time.Duration(timeout) * time.Second},
		"tcp"+strconv.Itoa(protocol),
		url+":"+strconv.Itoa(port),
		confTLS,
	)
	if err != nil {
		if ShowErrors {
			log.Printf(
				"%s: %s %s",
				Blue(url+":"+strconv.Itoa(port)),
				Red("Error during dial:"),
				Orange(err),
			)
		}
		return nil
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	return certs
}
