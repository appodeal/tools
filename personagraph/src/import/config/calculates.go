package config

import (
	"strings"
	"fmt"
	"os"

	"strconv"
)

type Calculate struct {
	Name        string
	CategoryIDs []int64
}

func (self *Calculate) String() string {
	return fmt.Sprintf("%s:%+v", self.Name, self.CategoryIDs)
}

type Calculates []Calculate

func (self *Calculates) String() string {
	calculates := make([]string, 0)

	for _, calculate := range *self {
		calculates = append(calculates, calculate.String())
	}

	return strings.Join(calculates, ", ")
}

func (self *Calculates) Set(value string) error {
	calculate := Calculate{}

	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		fmt.Printf("Invalid category calculate: %s\n", value)
		os.Exit(1)
	}
	calculate.Name = parts[0]


	for _, v := range strings.Split(parts[1],",") {
		if id, err := strconv.Atoi(v); err != nil {
			fmt.Printf("Invalid category id: %s\n", v)
			os.Exit(1)
		} else {
			calculate.CategoryIDs = append(calculate.CategoryIDs, int64(id))
		}
	}

	*self = append(*self, calculate)
	return nil
}
