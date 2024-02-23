package checkers

import (
	"fmt"

	chrome "github.com/MustafaNafizDurukan/CyberIndividualDefender/internal/checkers/browsers"
	"github.com/MustafaNafizDurukan/CyberIndividualDefender/pkg/types"
)

var allCheckerList = []types.IChecker{
	&chrome.Checker{},
}

func List() []types.IChecker {
	return allCheckerList
}

func Select(shortFlag string, longFlag string) types.IChecker {
	for _, checker := range allCheckerList {
		if checker.Descriptor().ShortFlag == shortFlag {
			return checker
		}

		if checker.Descriptor().LongFlag == longFlag {
			return checker
		}
	}

	return nil
}

func Init(checker types.IChecker) {
	checker.Init()
}

func Check(checker types.IChecker) bool {
	err := checker.Check()
	if err != nil {
		fmt.Printf("An error ocurred while checking with %s error: %v", checker.Descriptor().Name, err)
		return false
	}

	return true
}
