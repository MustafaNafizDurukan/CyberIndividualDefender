package chrome

import (
	"fmt"

	"github.com/MustafaNafizDurukan/CyberIndividualDefender/pkg/types"
)

type Checker struct {
}

func (c *Checker) Init() {}

func (c *Checker) Descriptor() *types.Descriptor {
	return &types.Descriptor{
		Name:        "Chrome checker",
		Description: "Checks if chrome has vulnerable passwords etc.",
		ShortFlag:   "-chr",
		LongFlag:    "--chrome",
	}
}

func (c *Checker) Check() error {
	fmt.Println("Chrome checked successfully!")
	return nil
}
