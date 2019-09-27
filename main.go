package main

import (
	"github.com/segmentio/go-log"
	"github.com/thechriswalker/system-stats/pkg/collector"
	"github.com/thechriswalker/system-stats/pkg/cpu"
	"github.com/thechriswalker/system-stats/pkg/disk"
	"github.com/thechriswalker/system-stats/pkg/memory"

	"os"
	"time"

	statsd "github.com/statsd/client"
	namespace "github.com/statsd/client-namespace"
	. "github.com/tj/go-gracefully"

	"github.com/tj/docopt"
)

// Version of this app
const Version = "0.3.0"

// Usage information
const Usage = `
  Usage:
    system-stats
      [--statsd-address addr]
      [--memory-interval i]
      [--memory-extended]
      [--disk-interval i]
      [--cpu-interval i]
      [--cpu-extended]
      [--name name]
    system-stats -h | --help
    system-stats --version

  Options:
    --statsd-address addr   statsd address [default: :8125]
    --memory-interval i     memory reporting interval [default: 10s]
    --memory-extended       output additional extended memory metrics
    --disk-interval i       disk reporting interval [default: 30s]
    --cpu-interval i        cpu reporting interval [default: 5s]
    --cpu-extended          output additional extended CPU metrics
    --name name             node name defaulting to hostname [default: hostname]
    -h, --help              output help information
    -v, --version           output version
`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	log.Check(err)

	log.Info("starting system %s", Version)

	client, err := statsd.Dial(args["--statsd-address"].(string))
	log.Check(err)

	name := args["--name"].(string)
	if "hostname" == name {
		host, err := os.Hostname()
		log.Check(err)
		name = host
	}

	c := collector.New(namespace.New(client, name))
	collectorCount := 0

	interval := getInterval(args, "--memory-interval")
	if interval > 0 {
		c.Add(memory.New(interval, args["--memory-extended"].(bool)))
		collectorCount++
	}
	interval = getInterval(args, "--cpu-interval")
	if interval > 0 {
		c.Add(cpu.New(interval, args["--cpu-extended"].(bool)))
		collectorCount++
	}
	interval = getInterval(args, "--disk-interval")
	if interval > 0 {
		c.Add(disk.New(interval))
		collectorCount++
	}

	c.Start()
	Shutdown()
	c.Stop()
}

func getInterval(args map[string]interface{}, name string) time.Duration {
	d, err := time.ParseDuration(args[name].(string))
	log.Check(err)
	return d
}
