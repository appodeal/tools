package application

import (
	"import/profile"
	"github.com/sirupsen/logrus"
	"github.com/aerospike/aerospike-client-go"
	"time"
	"sort"
)

type byWeightItem struct {
	Category int64
	Weight   float64
}

type byWeight []byWeightItem

func NewByWeight(categories []int64, weights []float64) []byWeightItem {
	items := make([]byWeightItem, 0)
	for index, category := range categories {
		items = append(items, byWeightItem{
			Category: category,
			Weight:   weights[index],
		})
	}
	return items
}

func (self byWeight) Len() int {
	return len(self)
}

func (self byWeight) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self byWeight) Less(i, j int) bool {
	return self[i].Weight > self[j].Weight
}

func (self *Application) Store(profile *profile.Profile, logger *logrus.Entry) error {
	var (
		key   *aerospike.Key
		table = self.Config.Table
	)

	if k, err := aerospike.NewKey(table.Namespace, table.Set, profile.Device); err != nil {
		return err
	} else {
		key = k
	}

	categories := NewByWeight(profile.Categories, profile.Weights)
	sort.Sort(byWeight(categories))

	ids := make([]int64, 0)
	for _, category := range categories {
		ids = append(ids, category.Category)
	}

	logger.Infof("%+v", categories)

	bins := aerospike.BinMap{
		"pg_updated_at": time.Now().Unix(),
	}

	if categories, err := self.Categories.ByIDs(ids...); err != nil {
		return err
	} else {
		bins["pg_segments"] = categories
	}

	if err := self.Aerospike.Put(nil, key, bins); err != nil {
		return err
	}

	return nil
}
