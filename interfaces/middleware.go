package interfaces

import (
	"context"
	"log"
	"net/http"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type Middleware struct {
	routes []*route
}

func (h *Middleware) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *Middleware) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			log.Printf("request from %s: path: %s", r.RemoteAddr, r.URL.Path)
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Content-Type", "application/json")
			if matchURL := route.pattern.FindStringSubmatch(r.URL.Path); len(matchURL) > 1 {
				r = r.WithContext(context.WithValue(r.Context(), "match", matchURL[1]))
			}
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
