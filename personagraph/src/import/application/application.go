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
	Config    *config.Config
	Logger    *logrus.Entry
	Importers chan string
	Group     sync.WaitGroup
	Aerospike *aerospike.Client
}

func New(config *config.Config, logger *logrus.Entry) (*Application, error) {
	application := &Application{Config: config, Logger: logger}
	application.Logger.Infof("Aerospike hosts: %s", config.Hosts.String())
	if client, err := aerospike.NewClientWithPolicyAndHost(aerospike.NewClientPolicy(), config.Hosts.Aerospike()...); err != nil {
		logger.Error(err)
		os.Exit(1)
	} else {
		application.Aerospike = client
	}
	application.Logger.Infof("Importers: %d", config.Importers)
	application.Importers = make(chan string, config.Importers)
	return application, nil
}

func (self *Application) Run() error {
	self.Group.Add(len(self.Config.Files))
	go func(self *Application) {
		for {
			file := <-self.Importers
			go func(self *Application, file string) {
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
