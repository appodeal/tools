package application

import (
	"import/config"
	"github.com/sirupsen/logrus"
	"github.com/aerospike/aerospike-client-go"
	"bufio"
	"sync"
	"os"
)

type Application struct {
	Config          *config.Config
	Logger          *logrus.Entry
	Importers       chan string
	Group           sync.WaitGroup
	Aerospike       *aerospike.Client
	Categories      Categories
	SkippedProfiles *os.File
}

func New(config *config.Config, logger *logrus.Entry) (*Application, error) {
	self := &Application{Config: config, Logger: logger}
	if !config.BlackHole {
		self.Logger.Infof("Connect to aerospike: %s", config.Hosts.String())
		if client, err := aerospike.NewClientWithPolicyAndHost(aerospike.NewClientPolicy(), config.Hosts.Aerospike()...); err != nil {
			self.Logger.Error(err)
			os.Exit(1)
		} else {
			self.Logger = self.Logger.WithField("storage", "aerospike")
			self.Aerospike = client
		}
	} else {
		self.Logger = self.Logger.WithField("storage", "black-hole")
	}

	if config.SkippedProfiles != "" {
		if file, err := os.OpenFile(config.SkippedProfiles, os.O_TRUNC | os.O_CREATE | os.O_WRONLY | os.O_EXCL, 0777); err != nil {
			logger.Error(err)
			os.Exit(1)
		} else {
			self.SkippedProfiles = file
		}
	}

	self.Logger.Infof("Loading categories from %s", config.Categories)
	if err := self.Categories.Load(config.Categories); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	self.Importers = make(chan string, config.Importers)
	return self, nil
}

func (self *Application) Run() error {
	// TODO: need fix parallel importing
	self.Group.Add(len(self.Config.Files))
	go func(self *Application) {
		for {
			file := <-self.Importers
			func(self *Application, file string) {
				defer self.Group.Done()
				self.Import(file, self.Logger.WithField("file", file))
			}(self, file)
		}
	}(self)

	for _, file := range self.Config.Files {
		self.Importers <- file
		self.Logger.WithField("file", file).Infof("Enqueued")
	}

	bufio.NewWriter(self.Logger.Logger.Out).Flush()
	self.Group.Wait()
	return nil
}
