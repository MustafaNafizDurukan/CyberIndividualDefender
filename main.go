package main

import (
	"fmt"

	"github.com/MustafaNafizDurukan/CyberIndividualDefender/internal/checkers"
)

func main() {
	checker := checkers.Select("-chr", "")
	if checker == nil {
		fmt.Println("The checker you requested was not found.")
		return
	}

	checkers.Init(checker)

	checkers.Check(checker)
}
