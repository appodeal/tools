package application

import (
	"github.com/sirupsen/logrus"
	"time"
	"os"
	"compress/gzip"
	"io"
	"bufio"
	"encoding/json"
	"import/profile"
)

func (self *Application) Import(path string, logger *logrus.Entry) error {
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
		var p profile.Profile
		Profiles++
		if err := json.Unmarshal([]byte(scanner.Text()), &p); err != nil {
			Errors++
			logger.Error(err)
		} else {
			Imported++
		}
	}
	return nil
}
