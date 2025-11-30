package globals

import (
	"fmt"
	"math/rand"
	"github.com/cespare/xxhash"
)

func Rand64() uint64 {
	return uint64(rand.Uint32()) << 32 | uint64(rand.Uint32())
}


func GetExchangeKey(s1, s2 string) uint64 {
	return xxhash.Sum64([]byte(GetExchangeName(s1,s2)))
}

func GetExchangeName(s1, s2 string) string {
	return fmt.Sprintf("%v:e:%v", s1, s2)
}