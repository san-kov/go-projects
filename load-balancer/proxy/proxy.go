package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Forward(w http.ResponseWriter, r *http.Request, targetURL string) {
	target, err := url.Parse("http://" + targetURL)

	if err != nil {
		http.Error(w, "invalid backend", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}
