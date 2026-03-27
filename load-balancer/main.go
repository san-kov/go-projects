package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL         string
	Alive       bool
	Connections int64
}

type LoadBalancer struct {
	backends []*Backend
	mu       sync.RWMutex
	counter  uint64
}

func (lb *LoadBalancer) AddBackend(url string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.backends = append(lb.backends, &Backend{URL: url, Alive: true})
}

func (lb *LoadBalancer) NextRoundRobin() (*Backend, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	total := uint64(len(lb.backends))

	for range lb.backends {
		idx := atomic.AddUint64(&lb.counter, 1)
		b := lb.backends[idx%total]

		if b.Alive {
			return b, nil
		}

	}

	return nil, fmt.Errorf("no backends available") // все мёртвые
}

func (lb *LoadBalancer) SetAlive(url string, alive bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	for _, b := range lb.backends {
		if b.URL == url {
			b.Alive = alive
		}
	}
}

func checkBackend(ctx context.Context, backend string) (string, error) {
	checkCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	delay := time.Duration(rand.Intn(300)) * time.Millisecond
	select {
	case <-time.After(delay):
		return backend + ": ok", nil
	case <-checkCtx.Done():
		return backend + ": timeout", checkCtx.Err()
	}

}

func checkBackends(ctx context.Context, backends []string) <-chan string {
	out := make(chan string)

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				for _, s := range backends {
					go func(s string) {
						result, err := checkBackend(ctx, s)

						if err != nil {
							return
						}
						select {
						case out <- result:
						case <-ctx.Done():
						}
					}(s)
				}
			case <-ctx.Done():
				close(out)
				return
			}
		}
	}()

	return out
}

func main() {
	// backends := []string{"backend-1:9000", "backend-2:9001", "backend-3:9002"}
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// results := checkBackends(ctx, backends)
	// for msg := range results {
	// 	fmt.Println(msg)
	// }
	// fmt.Println("health checks stopped")

	backend1 := Backend{URL: "backend1"}
	backend2 := Backend{URL: "backend2"}
	backend3 := Backend{URL: "backend3"}

	lb := LoadBalancer{backends: []*Backend{&backend1, &backend2, &backend3}}
	lb.SetAlive(backend1.URL, false)
	var wg sync.WaitGroup

	for i := 0; i < 9; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			backend, err := lb.NextRoundRobin()

			if err != nil {
				return
			}

			fmt.Println(backend.URL)
		}()
	}

	wg.Wait()
}
