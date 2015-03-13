package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type indexHandler struct {
	chain http.Handler
}

//IndexHandler redirects requests with no path to the root of Prefix
func IndexHandler(h http.Handler) http.Handler {
	return indexHandler{h}
}

func (h indexHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.RequestURI == config.Prefix {
		http.Redirect(rw, r, config.Prefix+"/", 301)
		return
	}
	h.chain.ServeHTTP(rw, r)
}

type forwardedHandler struct {
	chain http.Handler
}

//ForwardedHandler replaces the Remote Address with the X-Forwarded-For header if it exists
func ForwardedHandler(h http.Handler) http.Handler {
	return forwardedHandler{h}
}

func (h forwardedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	addr := strings.Split(r.RemoteAddr, ":")
	if len(addr) != 2 {
		log.Panicln("No Remote Address set")
	}

	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		r.RemoteAddr = fmt.Sprintf("%s:%s", ip, addr)
	}

	h.chain.ServeHTTP(rw, r)
}
