package service

import (
	"testing"
	"time"

	"github.com/cyanly/gotrade/core/messagebus/test"
	proto "github.com/cyanly/gotrade/proto/service"
)

func TestServiceStartAndStop(t *testing.T) {
	ts := test.RunDefaultServer()
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
