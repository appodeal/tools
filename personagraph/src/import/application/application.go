package application

import (
	"import/config"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config *config.Config
	Logger *logrus.Entry
}

func New(config *config.Config, logger *logrus.Entry) (*Application, error) {
	application := &Application{Config: config, Logger: logger}
	return application, nil
}

func (self *Application) Run() error {
	for _, name := range self.Config.Files {
		if err := self.Import(name, self.Logger.WithField("file", name)); err != nil {
			return err
		}
	}
	return nil
}
