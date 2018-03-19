package config

import (
	"strings"
	"fmt"
	"os"

	"strconv"
)

type Filter struct {
	Name        string
	CategoryIDs []int64
}

func (self *Filter) String() string {
	return fmt.Sprintf("%s:%+v", self.Name, self.CategoryIDs)
}

type Filters []Filter

func (self *Filters) String() string {
	filters := make([]string, 0)

	for _, filter := range *self {
		filters = append(filters, filter.String())
	}

	return strings.Join(filters, ", ")
}

func (self *Filters) NameByIDs(ids ... int64) []string {
	names := make([]string, 0)
	for _, filter := range *self {
		for _, v1 := range ids {
			for _, v2 := range filter.CategoryIDs {
				if v1 == v2 {
					names = append(names, filter.Name)
				}
			}
		}
	}
	return names
}

func (self *Filters) Set(value string) error {
	filter := Filter{}

	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		fmt.Printf("Invalid category calculate: %s\n", value)
		os.Exit(1)
	}
	filter.Name = parts[0]

	for _, v := range strings.Split(parts[1], ",") {
		if id, err := strconv.Atoi(v); err != nil {
			fmt.Printf("Invalid category id: %s\n", v)
			os.Exit(1)
		} else {
			filter.CategoryIDs = append(filter.CategoryIDs, int64(id))
		}
	}

	*self = append(*self, filter)
	return nil
}
