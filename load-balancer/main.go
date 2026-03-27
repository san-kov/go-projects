package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

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
	backends := []string{"backend-1:9000", "backend-2:9001", "backend-3:9002"}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results := checkBackends(ctx, backends)
	for msg := range results {
		fmt.Println(msg)
	}
	fmt.Println("health checks stopped")

}
