package cache_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/SoliMark/gotasker-pro/internal/cache"
)

func TestJitterTTLRange(t *testing.T) {
	src := rand.New(rand.NewSource(1))
	j := cache.Jitter{Rand: src}
	base := 60 * time.Second
	lo, hi := base-base/10, base+base/10

	for i := 0; i < 1000; i++ {
		got := j.TTL(base, 0.1, time.Second)
		if got < lo || got > hi {
			t.Fatalf("out of range: %v", got)
		}
	}
}
