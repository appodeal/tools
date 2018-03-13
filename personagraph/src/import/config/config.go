package config

import (
	"flag"
	"os"
	"fmt"
	"strings"
	"github.com/aerospike/aerospike-client-go"
	"net"
	"strconv"
)

type Host struct {
	Address string
	Port    int
}

func (self *Host) String() string {
	return fmt.Sprintf("%s:%d", self.Address, self.Port)
}

type Hosts []Host

func (self *Hosts) String() string {
	hosts := make([]string, 0)

	for _, host := range *self {
		hosts = append(hosts, host.String())
	}

	return strings.Join(hosts, ", ")
}

func (self *Hosts) Aerospike() []*aerospike.Host {
	hosts := make([]*aerospike.Host, 0)

	for _, host := range *self {
		hosts = append(hosts, &aerospike.Host{
			Name: host.Address,
			Port: host.Port,
		})
	}
	return hosts
}

func (self *Hosts) Set(value string) error {
	address, port, err := net.SplitHostPort(value)

	if err != nil {
		fmt.Printf("Invalid aerospike address(%v): %s\n", err, value)
		os.Exit(1)
	}

	host := Host{Address: address}

	if port == "" {
		host.Port = 3000
	} else if p, err := strconv.Atoi(port); err != nil {
		fmt.Printf("Invalid aerospike port(%v): %s\n", err, port)
		os.Exit(1)
	} else {
		host.Port = p
	}

	*self = append(*self, host)
	return nil
}

type Table struct {
	Namespace string
	Set string
}

type Config struct {
	Files     []string
	Remove    bool
	Hosts     Hosts
	Importers int
	Table	 Table
}

func New() (*Config, error) {
	config := &Config{
		Remove:    false,
		Hosts:     Hosts{},
		Importers: 1,
		Table: Table{
			Namespace: "appodeal",
			Set: "device",
		},
	}

	options := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] <files>\n", os.Args[0])
		options.PrintDefaults()
	}

	options.Var(&config.Hosts, "s", "Aerospike hosts")
	options.StringVar(&config.Table.Namespace, "n", config.Table.Namespace, "Aerospike namespace")
	options.StringVar(&config.Table.Set, "t", config.Table.Set, "Aerospike set")
	options.BoolVar(&config.Remove, "r", config.Remove, "remove file after processing")
	options.IntVar(&config.Importers, "i", config.Importers, "number of parallel importers")
	options.Parse(os.Args[1:])

	if len(config.Hosts) == 0 {
		config.Hosts = Hosts{Host{"127.0.0.1", 3000}}
	}

	config.Files = options.Args()

	return config, nil
}
