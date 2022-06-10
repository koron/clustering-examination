package main

import (
	"fmt"

	"gonum.org/v1/gonum/stat/combin"
)

func main() {
	fmt.Printf("P(20, 20)=%d\n", combin.NumPermutations(20, 20))
}
