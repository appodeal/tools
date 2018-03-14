package config

import (
	"fmt"
	"strings"
	"github.com/aerospike/aerospike-client-go"
	"net"
	"os"
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
