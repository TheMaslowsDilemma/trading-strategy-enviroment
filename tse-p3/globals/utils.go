package globals

import (
	"math/rand"
)

func Rand64() uint64 {
	return uint64(rand.Uint32()) << 32 | uint64(rand.Uint32())
}
