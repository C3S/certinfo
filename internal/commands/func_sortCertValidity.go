package commands

import (
	"sort"

	. "github.com/C3S/certinfo/internal/globals"
)

func sortCertValidity(certs map[string]*CertValidity) []string {
	keys := make([]string, 0, len(certs))

	for key := range certs {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		var_i := *certs[keys[i]]
		var_j := *certs[keys[j]]
		return var_i.DaysLeft < var_j.DaysLeft
	})

	return keys
}
