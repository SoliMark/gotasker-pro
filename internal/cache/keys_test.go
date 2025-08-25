package cache_test

import (
	"testing"

	"github.com/SoliMark/gotasker-pro/internal/cache"
)

func TestKeyUserTasks(t *testing.T) {
	k := cache.KeyUserTasks(42)
	if k != "user:42:tasks:v1" {
		t.Fatalf("got %q", k)
	}
}
