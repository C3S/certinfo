package globals

import "crypto/x509"

type Host struct {
	URL  string `mapstructure:"url"`
	Port int    `mapstructure:"port"`
}

type OddEvenKeys struct {
	Odd  []string
	Even []string
}

type CertValidity struct {
	URL         string
	Port        int
	Protocol    int
	Certificate []*x509.Certificate
	DaysLeft    int
}
