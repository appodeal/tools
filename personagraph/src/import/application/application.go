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
	Queue           chan string
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
	self.Queue = make(chan string)
	return self, nil
}

func (self *Application) Run() error {
	self.Group.Add(len(self.Config.Files))

	for i := 0; i < self.Config.Importers; i++ {
		go func(self *Application, queue chan string) {
			defer self.Group.Done()
			file := <- queue
			self.Import(file, self.Logger.WithField("file", file))
		}(self, self.Queue)
	}

	for _, file := range self.Config.Files {
		self.Logger.WithField("file", file).Infof("Enqueued")
		self.Queue <- file
	}

	bufio.NewWriter(self.Logger.Logger.Out).Flush()
	self.Group.Wait()
	return nil
}
