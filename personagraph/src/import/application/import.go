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
	"path/filepath"
	"github.com/aerospike/aerospike-client-go"
	"github.com/aerospike/aerospike-client-go/types"
	"fmt"
)

func (self *Application) Import(path string, logger *logrus.Entry) error {
	StartAt := time.Now()
	Imported := 0
	Total := 0
	calculates := map[string]int64{}
	defer (func() {
		stats := logrus.Fields{}
		for key, value := range calculates {
			stats[key] = value
		}
		logger.WithFields(logrus.Fields{
			"imported": Imported,
			"total": Total,
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

	policy := aerospike.WritePolicy{
		BasePolicy:         *aerospike.NewPolicy(),
		GenerationPolicy:   aerospike.NONE,
		CommitLevel:        aerospike.COMMIT_ALL,
		Generation:         0,
		Expiration:         0,
		SendKey:            false,
	}

	if self.Config.UpdateOnly {
		policy.RecordExistsAction = aerospike.UPDATE_ONLY
	} else {
		policy.RecordExistsAction = aerospike.UPDATE
	}

	scanner := bufio.NewScanner(zip)
	for scanner.Scan() {
		var p profile.Profile
		text := scanner.Text()
		if err := json.Unmarshal([]byte(text), &p); err != nil {
			logger.Error(err)
			return err
		} else {
			Total++
			if len(self.Config.Calculates) > 0 {
				if name := self.Config.Calculates.NameByIDs(p.Categories...); name != "" {
					if v, ok := calculates[name]; ok {
						calculates[name] = v + 1
					} else {
						calculates[name] = 1
					}
				} else if self.SkippedProfiles != nil {
					fmt.Fprintln(self.SkippedProfiles, text)
					continue
				}
			}
			if !self.Config.BlackHole {
				if err := self.Store(&policy, &p, logger.WithField("source", "store")); err != nil {
					ae := err.(types.AerospikeError)
					if !(self.Config.UpdateOnly && ae.ResultCode() == types.KEY_NOT_FOUND_ERROR) {
						logger.WithField("code", ae.ResultCode()).Error(ae)
						return err
					}
				}
			}
			Imported++
		}
	}
	if moveTo := self.Config.MoveTo; moveTo != "" {
		if err := os.Rename(path, filepath.Join(moveTo, filepath.Base(file.Name()))); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	}
	return nil
}
