package tools

import "math/rand"

func SourceFromSeed(seed string) int64 {
	ss := int64(1323123)
	for i := 0; i < len(seed); i++ {
		ss += int64(seed[i]) + 1
	}
	return ss
}

func RandFromSeed(seed int64) *rand.Rand {
	source := rand.NewSource(seed)
	return rand.New(source)
}
