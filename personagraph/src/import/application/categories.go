package application

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
)

type Categories map[int64]string

func (self *Categories) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	if err := yaml.NewDecoder(file).Decode(self); err != nil {
		return err
	}
	return nil
}

func (self *Categories) ByIDs(ids ... int64) ([]string, error) {
	values := make([]string, 0)
	for _, id := range ids {
		if value, ok := (*self)[id]; !ok {
			return []string{}, errors.New(fmt.Sprintf("Unknown category id: %d", id))
		} else {
			values = append(values, value)
		}
	}
	return values, nil
}
