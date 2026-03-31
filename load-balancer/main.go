package main

import (
	"context"
	"fmt"
	"load-balancer/balancer"
	"net/http"
	"time"
)

func startBackend(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "response from backend %s\n", port)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	go http.ListenAndServe(":"+port, mux)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startBackend("9001")
	startBackend("9002")
	startBackend("9003")
	time.Sleep(100 * time.Millisecond)

	lb := &balancer.LoadBalancer{}
	lb.AddBackend("localhost:9001")
	lb.AddBackend("localhost:9002")
	lb.AddBackend("localhost:9003")

	lb.StartHealthChecks(ctx)

	fmt.Println("Load balancer listening on :8080")
	http.ListenAndServe(":8080", lb)

	if err := http.ListenAndServe(":8080", lb); err != nil {
		fmt.Println("server error:", err)
	}
}
