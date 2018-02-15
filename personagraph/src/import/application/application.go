package application

import (
	"import/config"
	"github.com/sirupsen/logrus"
	"os"
	"compress/gzip"
	"io"
	"bufio"
	"encoding/json"
	"time"
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
		if err := self.ProcessingFile(name, self.Logger.WithField("file", name)); err != nil {
			return err
		}
	}
	return nil
}

type Profile struct {
	ID         int64     `json:"u"`
	Categories []int64   `json:"a"`
	Weights    []float64 `json:"w"`
	Device     string    `json:"d"`
}

func (self *Application) ProcessingFile(path string, logger *logrus.Entry) error {
	StartAt := time.Now()
	Profiles := 0
	Imported := 0
	Errors := 0
	defer (func() {
		logger.WithFields(logrus.Fields{
			"imported": Imported,
			"errors": Errors,
			"profiles": Profiles,
			"elapsed": time.Since(StartAt),
		}).Info("Finish import")
	})()
	logger.Infof("Start import")
	file, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer file.Close()
	zip, err := gzip.NewReader(file)
	if err == io.EOF { return nil }
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer zip.Close()
	scanner := bufio.NewScanner(zip)
	for scanner.Scan() {
		var profile Profile
		Profiles++
		if err := json.Unmarshal([]byte(scanner.Text()), &profile); err != nil {
			Errors++
			logger.Error(err)
		} else {
			Imported++
		}
	}
	return nil
}
