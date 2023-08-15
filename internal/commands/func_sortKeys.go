package commands

import (
	"sort"

	. "github.com/C3S/certinfo/internal/globals"
)

func sortKeys(hosts map[string]*Host) []string {
	keys := make([]string, 0, len(hosts))

	for k := range hosts {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
