package random

import (
	crand "crypto/rand"
	"encoding/base64"
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

// GenerateRandomBytes returns securely generated random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded securely generated random string
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
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
