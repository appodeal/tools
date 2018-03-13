package application

import (
	"import/profile"
	"github.com/sirupsen/logrus"
	"github.com/aerospike/aerospike-client-go"
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
		"test": 3,
	}

	if err := self.Aerospike.Put(nil, key, bins); err != nil {
		return err
	}

	return nil
}