package dial

import "time"

func DialTimeout(sec int) time.Duration {
	return time.Duration(sec) * time.Second
}
