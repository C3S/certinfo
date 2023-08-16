package commands

import (
	"fmt"
	"strconv"

	. "github.com/C3S/certinfo/internal/globals"
)

func bestBeforeCheck(
	certs map[string]*CertValidity,
) {
	sortedKeys := sortCertValidity(certs)
	for _, i := range sortedKeys {
		if certs[i].DaysLeft > Days {
			fmt.Printf(
				"%-35s (IPv%d): expires %-44s %s",
				Blue(certs[i].URL),
				certs[i].Protocol,
				Green(certs[i].Certificate[0].NotAfter.Format("02.01.2006"))+", in "+Green(strconv.Itoa(certs[i].DaysLeft))+" days",
				Green("-- ok!\n"),
			)
			continue
		} else if certs[i].DaysLeft < 0 {
			fmt.Printf(
				"%-35s (IPv%d): expired %-44s %s",
				Blue(certs[i].URL),
				certs[i].Protocol,
				Red(certs[i].Certificate[0].NotAfter.Format("02.01.2006"))+", "+Red(strconv.Itoa(certs[i].DaysLeft))+" days ago",
				Red("-- red alert!\n"),
			)
			continue
		} else {
			fmt.Printf(
				"%-35s (IPv%d): expires %-44s %s",
				Blue(certs[i].URL),
				certs[i].Protocol,
				Orange(certs[i].Certificate[0].NotAfter.Format("02.01.2006"))+", in "+Orange(strconv.Itoa(certs[i].DaysLeft))+" days",
				Orange("-- please renew!\n"),
			)
			continue
		}
	}
}
