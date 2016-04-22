package test

import (
	"github.com/nats-io/gnatsd/server"
	"github.com/nats-io/gnatsd/test"
)

// Create a local messaging broker, usually for testing purpose
func RunDefaultServer() *server.Server {
	test.DefaultTestOptions.Port = 22222
	return test.RunDefaultServer()
}
