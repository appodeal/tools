package config

import (
	"flag"
	"os"
	"fmt"
)

type Table struct {
	Namespace string
	Set       string
}

type Config struct {
	Files      []string
	Hosts      Hosts
	Importers  int
	Table      Table
	Categories string
	Calculates Calculates
	BlackHole  bool
	MoveTo     string
}

func New() (*Config, error) {
	config := &Config{
		Hosts:     Hosts{},
		Importers: 1,
		Table: Table{
			Namespace: "appodeal",
			Set:       "device_apps_interests",
		},
		Categories: "./categories.yml",
	}

	options := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] <files>\n", os.Args[0])
		options.PrintDefaults()
	}

	options.Var(&config.Hosts, "s", "Aerospike hosts (example 127.0.0.1:3000)")
	options.StringVar(&config.Table.Namespace, "n", config.Table.Namespace, "Aerospike namespace")
	options.StringVar(&config.Table.Set, "t", config.Table.Set, "Aerospike set")
	options.StringVar(&config.Categories, "p", config.Categories, "YAML file with categories from Personagraph")
	options.IntVar(&config.Importers, "i", config.Importers, "number of parallel importers")
	options.Var(&config.Calculates, "c", "calculate of category (example: name:1,2,3,4)")
	options.BoolVar(&config.BlackHole, "b", config.BlackHole, "don't write profiles to aerospike")
	options.StringVar(&config.MoveTo, "m", config.MoveTo, "move dump files after import to directory")
	options.Parse(os.Args[1:])

	if len(config.Hosts) == 0 {
		config.Hosts = Hosts{Host{"127.0.0.1", 3000}}
	}

	config.Files = options.Args()

	return config, nil
}
