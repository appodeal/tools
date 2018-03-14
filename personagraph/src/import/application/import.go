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
	Imported := 0
	calculates := map[string]int64{}
	defer (func() {
		stats := logrus.Fields{}
		for key, value := range calculates {
			stats[key] = value
		}
		logger.WithFields(logrus.Fields{
			"imported": Imported,
			"elapsed":  time.Since(StartAt),
		}).WithFields(stats).Info("Finish import")
	})()
	logger.Infof("Start import")
	file, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer file.Close()
	zip, err := gzip.NewReader(file)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer zip.Close()
	scanner := bufio.NewScanner(zip)
	for scanner.Scan() {
		var p profile.Profile
		if err := json.Unmarshal([]byte(scanner.Text()), &p); err != nil {
			logger.Error(err)
			return err
		} else {
			if name := self.Config.Calculates.NameByIDs(p.Categories...); name != "" {
				if v, ok := calculates[name]; ok {
					calculates[name] = v + 1
				} else {
					calculates[name] = 1
				}
			}
			if self.Config.BlackHole {
				continue
				Imported++
			} else {
				if err := self.Store(&p, logger.WithField("source", "store")); err != nil {
					logger.Error(err)
					return err
				} else {
					Imported++
				}
			}
		}
	}
	return nil
}
