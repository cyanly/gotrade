package messagebus

import "github.com/nats-io/nats.go"

// Config represents the configuration.
type Config struct {
	Url string
}

const DefaultUrl = nats.DefaultURL

// NewConfig returns an instance of Config with defaults.
func NewConfig(url string) *Config {

	var config = &Config{
		Url: url,
	}
	return config
}
