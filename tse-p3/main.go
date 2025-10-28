package main

import (
	"fmt"
	"github.com/holiman/uint256"
)

func main() {
	fmt.Println("--- Trading Stategy Environment: Part Three ---")
	
	var (
		c1    *uint256.Int
		c2    *uint256.Int
		big1  *uint256.Int
		big2  *uint256.Int
		bigres *uint256.Int
	)

	c1    = uint256.NewInt(7829)
	c2    = uint256.NewInt(6703)
	big1  = uint256.NewInt(123456)
	big2  = uint256.NewInt(987654)
	bigres = uint256.NewInt(0)

	for i := 0; i < 40; i++ {
		bigres.Add(big1, big2)
		fmt.Printf("%02d: %v + %v = %v\n", i, big1, big2, bigres)
		big1.Add(big1, c1)
		big2.Or(big2.Add(big2, c2), bigres)
	}


}
