package getSnapshot

import (
	"sync"
	"testing"
	"time"
)

func TestGetSnapshot(t *testing.T) {
	start := time.Now()
	t.Log(GetSnapshot("https://cache.icpunks.com/icpunks/Token/%d", 1, 100))
	t.Log(time.Now().Sub(start))
}

func TestGetSnapshotPara(t *testing.T) {
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		from := i * 20
		to := i*20 + 20
		if i == 0 {
			from = 1
		}
		go func(from, to int) {
			wg.Add(1)
			defer wg.Done()
			t.Log("from: ", from, "to: ", to)
			t.Log(GetSnapshot("https://cache.icpunks.com/icpunks/Token/%d", from, to))
		}(from, to)
	}
	wg.Wait()
	t.Log(time.Now().Sub(start))
}
