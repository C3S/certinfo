package commands

import (
	"crypto/tls"
	"sync"
	"time"

	. "github.com/C3S/certinfo/internal/globals"
	"github.com/schollz/progressbar/v3"
)

func workerDaysValid(
	keys []string,
	now time.Time,
	jobs <-chan int,
	results chan<- map[string]*CertValidity,
	bar *progressbar.ProgressBar,
	confTLS *tls.Config,
) {
	for j := range jobs {
		thisKey := []string{keys[j]}
		results <- daysValid(thisKey, AllHosts, IPversions, Timeout, confTLS, now)
		bar.Add(1)
	}
}

func BestBeforeHosts(
	confTLS *tls.Config,
) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	keys := sortKeys(AllHosts)
	nKeys := len(keys)
	jobs := make(chan int, nKeys)
	results := make(chan map[string]*CertValidity, nKeys)
	bar := progressbar.NewOptions(
		len(keys),
		progressbar.OptionSetDescription("Checking host certificates..."),
	)

	nWorkers := 4
	var wg sync.WaitGroup
	for i := 1; i <= nWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerDaysValid(keys, now, jobs, results, bar, confTLS)
		}()
	}

	for j := 0; j < nKeys; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	for r := range results {
		bestBeforeCheck(r)
	}
}
