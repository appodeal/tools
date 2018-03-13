package application

import (
	"import/profile"
	"github.com/sirupsen/logrus"
	"github.com/aerospike/aerospike-client-go"
	"time"
)

func (self *Application) Store(profile *profile.Profile, logger *logrus.Entry) error {
	var (
		key *aerospike.Key
		table = self.Config.Table
	)

	if k, err := aerospike.NewKey(table.Namespace, table.Set, profile.Device); err != nil {
		return err
	} else {
		key = k
	}

	bins := aerospike.BinMap{
		"pg_updated_at": time.Now().Unix(),
	}

	if categories, err := self.Categories.ByIDs(profile.Categories...); err != nil {
		return err
	} else {
		bins["pg_segments"] = categories
	}

	if err := self.Aerospike.Put(nil, key, bins); err != nil {
		return err
	}

	return nil
}