package balancer

import (
	"fmt"
	"load-balancer/proxy"
	"net/http"
	"sync"
	"sync/atomic"
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

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend, err := lb.NextRoundRobin()

	if err != nil {
		http.Error(w, "no backends available", http.StatusServiceUnavailable)
		return
	}

	proxy.Forward(w, r, backend.URL)
}
