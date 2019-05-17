package random

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	rand "math/rand"
)

// GenerateRandomInt generated random int between min and max(inclusive)
func GenerateRandomInt(min int64, max int64) int64 {
	if max <= min {
		return min
	}

	var src cryptoSource
	rnd := rand.New(src)

	return rnd.Int63n((max+1)-min) + min
}

// Helper to turn crypto/rand into a seeder for math/rand
type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}
func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}
func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
