package service

import "os"

var (
	uuid     string
	hostname string
	pid      int
)

func init() {
	hn, _ := os.Hostname()
	hostname = hn

	pid = os.Getpid()
}
