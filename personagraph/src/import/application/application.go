package application

import (
	"import/config"
	"github.com/sirupsen/logrus"
	"github.com/aerospike/aerospike-client-go"
)

type Application struct {
	Config *config.Config
	Logger *logrus.Entry
}

func New(config *config.Config, logger *logrus.Entry) (*Application, error) {
	application := &Application{Config: config, Logger: logger}
	application.Logger.Infof("Aerospike hosts: %s", config.Hosts.String())
	aerospike.NewClientWithPolicyAndHost(aerospike.NewClientPolicy(), config.Hosts.Aerospike()...)
	return application, nil
}

func (self *Application) Run() error {
	for _, name := range self.Config.Files {
		self.Logger.Info("..")
		if err := self.Import(name, self.Logger.WithField("file", name)); err != nil {
			return err
		}
	}
	return nil
}
