package devserver

import (
	"fmt"
	"net/http"
	"sync"
)

type Reloader struct {
	mu      sync.Mutex
	clients map[chan string]struct{}
}

func NewReloader() *Reloader {
	return &Reloader{
		clients: make(
			map[chan string]struct{},
		),
	}
}

func (r *Reloader) Handler(
	w http.ResponseWriter,
	req *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"text/event-stream",
	)

	w.Header().Set(
		"Cache-Control",
		"no-cache",
	)

	ch := make(chan string)

	r.mu.Lock()
	r.clients[ch] = struct{}{}
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		delete(r.clients, ch)
		r.mu.Unlock()
	}()

	for msg := range ch {

		fmt.Fprintf(
			w,
			"data: %s\n\n",
			msg,
		)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func (r *Reloader) Broadcast() {

	r.mu.Lock()
	defer r.mu.Unlock()

	for ch := range r.clients {

		select {
		case ch <- "reload":
		default:
		}
	}
}
