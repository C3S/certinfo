package commands

import . "github.com/C3S/certinfo/internal/globals"

// split keys in two halves alternating
func oddEvenSplit(keys []string) OddEvenKeys {
	var oddKeys []string
	var evenKeys []string
	for x := 0; x <= len(keys)-1; x++ {
		if x%2 != 0 {
			oddKeys = append(oddKeys, keys[x])
		} else {
			evenKeys = append(evenKeys, keys[x])
		}
	}
	result := OddEvenKeys{
		Odd:  oddKeys,
		Even: evenKeys,
	}
	return (result)
}
