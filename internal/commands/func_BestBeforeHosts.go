package commands

import (
	"crypto/tls"
	"sync"
	"time"

	. "github.com/C3S/certinfo/internal/globals"
)

func BestBeforeHosts(confTLS *tls.Config) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	keys := sortKeys(AllHosts)
	splitKeys := oddEvenSplit(keys)
	splitKeysOdd := oddEvenSplit(splitKeys.Odd)
	splitKeysEven := oddEvenSplit(splitKeys.Even)
	wg := new(sync.WaitGroup)
	wg.Add(4)
	go func(wg *sync.WaitGroup) {
		bestBeforeCheck(splitKeysOdd.Odd, AllHosts, IPversions, Timeout, confTLS, now)
		wg.Done()
	}(wg)
	go func(wg *sync.WaitGroup) {
		bestBeforeCheck(splitKeysOdd.Even, AllHosts, IPversions, Timeout, confTLS, now)
		wg.Done()
	}(wg)
	go func(wg *sync.WaitGroup) {
		bestBeforeCheck(splitKeysEven.Odd, AllHosts, IPversions, Timeout, confTLS, now)
		wg.Done()
	}(wg)
	go func(wg *sync.WaitGroup) {
		bestBeforeCheck(splitKeysEven.Even, AllHosts, IPversions, Timeout, confTLS, now)
		wg.Done()
	}(wg)
	wg.Wait()
}
