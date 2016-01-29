package service

import (
	proto "github.com/cyanly/gotrade/proto/service"
	gnatsd "github.com/nats-io/gnatsd/test"
	"testing"
	"time"
)

func TestServiceStartAndStop(t *testing.T) {
	gnatsd.DefaultTestOptions.Port = 22222
	ts := gnatsd.RunDefaultServer()
	defer ts.Shutdown()

	sc := NewConfig()
	sc.ServiceName = "Test Service"
	sc.MessageBusURL = "nats://localhost:22222"

	svc := NewService(sc)
	if svc.Status != proto.STARTING {
		t.Fatalf("unexpected service status: %s , expecting STARTING", svc.Status)
	}

	svc.Start()
	if svc.Status != proto.RUNNING {
		t.Fatalf("unexpected service status: %s, expecting RUNNING", svc.Status)
	}

	time.Sleep(500 * time.Millisecond)
	svc.shutdownChannel <- true
	time.Sleep(100 * time.Millisecond)

	if svc.Status != proto.STOPPED {
		t.Fatalf("unexpected service status: %s, expecting STOPPED", svc.Status)
	}
}
