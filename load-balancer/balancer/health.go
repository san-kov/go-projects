package balancer

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (lb *LoadBalancer) check(ctx context.Context, url string) bool {
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(checkCtx, "GET", "http://"+url+"/health", nil)

	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (lb *LoadBalancer) checkAll(ctx context.Context) {
	lb.mu.RLock()
	backends := make([]*Backend, len(lb.backends))
	copy(backends, lb.backends)
	lb.mu.RUnlock()

	for _, b := range backends {
		go func(backend *Backend) {
			alive := lb.check(ctx, backend.URL)
			lb.SetAlive(backend.URL, alive)
			fmt.Printf("health check %s: alive=%v\n", backend.URL, alive)
		}(b)
	}

}

func (lb *LoadBalancer) StartHealthChecks(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				lb.checkAll(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}
