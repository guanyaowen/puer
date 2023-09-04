package util

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// RandomNumber 随机ID
func RandomNumber() string {
	seedGenerator := NewRand(time.Now().UnixNano())
	generator := rand.NewSource(seedGenerator.Int63())
	return fmt.Sprintf("%016x", uint64(generator.Int63()))
}

type lockedSource struct {
	mut sync.Mutex
	src rand.Source
}

func NewRand(seed int64) *rand.Rand {
	return rand.New(&lockedSource{src: rand.NewSource(seed)})
}

func (r *lockedSource) Int63() int64 {
	r.mut.Lock()
	defer r.mut.Unlock()
	return r.src.Int63()
}

// Seed implements Seed() of Source
func (r *lockedSource) Seed(seed int64) {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.src.Seed(seed)
}
