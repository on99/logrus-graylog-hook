package graylog

import (
	"compress/flate"
	"errors"
	"os"
	"time"

	"github.com/Graylog2/go-gelf/gelf"
)

// Config stores graylog hook config
type Config struct {
	Hostname   string
	Facility   string
	StaticMeta map[string]interface{} // static meta always sent to graylog

	// default not compress to save cpu power as logs are most likely to send over private network
	CompressType     gelf.CompressType // default to NoCompression
	CompressionLevel int               // default to flate.NoCompression

	// health check
	HealthCheckInterval time.Duration // default to 5 seconds
}

// NodeConfig stores node configuration
type NodeConfig struct {
	UDPAddress     string
	HealthCheckURL string
	Weight         int
}

// NewConfig returns a default graylog hook config
func NewConfig() Config {
	hostname, _ := os.Hostname()

	return Config{
		Hostname: hostname,

		CompressType:     gelf.CompressNone,
		CompressionLevel: flate.NoCompression,

		HealthCheckInterval: time.Second * 5,
	}
}

// ValidateConfig validates config
func ValidateConfig(c Config) error {
	if c.Hostname == "" {
		return errors.New("hostname cannot be empty")
	}

	return nil
}
