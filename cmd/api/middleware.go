package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// function that recovers from panics in our main application
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// checking if there has been a panic
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				// generate a 500 error from errors.go
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// function that limits requests to 2 per second with a max burst of 4,
// configured to limit by IP address
func (app *application) rateLimit(next http.Handler) http.Handler {

	// client struct to track last seen time
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// background goroutine to remove old entries from clients map,
	// to reduce resource use of clients map
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.config.limiter.enabled {
			// get user's IP address
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// Locks mutex to prevent race conditions
			mu.Lock()

			// checks for IP, then adds IP and limiter to map
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			// update last seen time for client
			clients[ip].lastSeen = time.Now()

			// if client exceeded limit, send a Too Many Request response
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			// Unlock mutex before calling next handler
			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
