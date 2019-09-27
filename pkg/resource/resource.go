package resource

import statsd "github.com/statsd/client-interface"

// Resource describes something to measure
type Resource interface {
	Name() string
	Start(statsd.Client) error
	Stop() error
}
