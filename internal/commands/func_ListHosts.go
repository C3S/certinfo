package commands

import (
	"fmt"
	"strconv"

	. "github.com/C3S/certinfo/internal/globals"
)

func ListHosts(hosts map[string]*Host) {
	keys := sortKeys(hosts)
	for _, i := range keys {
		fmt.Printf("%s: %s\n", Blue(hosts[i].URL), Orange(strconv.Itoa(hosts[i].Port)))
	}
}
