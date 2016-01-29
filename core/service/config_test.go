package service

import (
	"encoding/json"
	"testing"
)

func TestConfig_Parse(t *testing.T) {
	// Parse configuration.
	var c Config
	if err := json.Unmarshal([]byte(`{
"service_name": "test_service name",
"service_heartbeat_freq": 3,
"service_message_bus": "some_url",
"entry_not_in_use": "#####"
}`), &c); err != nil {
		t.Fatal(err)
	}

	// Validate configuration.
	if c.ServiceName != "test_service name" {
		t.Fatalf("unexpected service_name: %v", c.ServiceName)
	} else if c.HeartbeatFreq != 3 {
		t.Fatalf("unexpected service_heartbeat_freq: %v", c.HeartbeatFreq)
	}
}
