package main

import (
	"import/config"
	"import/application"
	"github.com/sirupsen/logrus"
	"formatter"
)

func init() {
}

func main() {
	logger := logrus.New()

	logger.Formatter = &formatter.TextFormatter{}
	if c, err := config.New(); err != nil {
		logger.Error(err)
	} else if a, err := application.New(c, logrus.NewEntry(logger)); err != nil {
		logger.Error(err)
	} else if err := a.Run(); err != nil {
		logger.Error(err)
	}
}
