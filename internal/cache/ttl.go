package cache

import (
	"math/rand"
	"time"
)

// Jitter 允許注入亂數源，方便測試；Rand 為 nil 時使用全域 rand。
type Jitter struct {
	Rand *rand.Rand
}

// TTL 以 base 為基準，加入 ±ratio 的隨機抖動。
// 會確保結果 > 0；若給 minTTL，則不小於 minTTL。
func (j Jitter) TTL(base time.Duration, ratio float64, minTTL ...time.Duration) time.Duration {
	if base <= 0 || ratio <= 0 {
		return base
	}
	if ratio > 1 {
		ratio = 1
	}
	r := rand.Float64
	if j.Rand != nil {
		r = j.Rand.Float64
	}
	delta := time.Duration((r()*2 - 1) * ratio * float64(base)) // ±ratio
	out := base + delta
	if len(minTTL) > 0 && out < minTTL[0] {
		out = minTTL[0]
	}
	if out <= 0 {
		out = time.Second // 保險：避免 0/負值導致立刻過期
	}
	return out
}
