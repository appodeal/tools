package config

import (
	"flag"
	"os"
	"fmt"
)

type Config struct {
	Files  []string
	Remove bool
}

func New() (*Config, error) {
	config := &Config{
		Remove: false,
	}

	options := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	options.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] <files>\n", os.Args[0])
		options.PrintDefaults()
	}
	options.BoolVar(&config.Remove, "r", config.Remove, "remove file after processing")
	options.Parse(os.Args[1:])
	config.Files = options.Args()

	return config, nil
}
