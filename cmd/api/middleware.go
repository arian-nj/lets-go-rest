package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/arian-nj/site/back/internal/data"
	"github.com/arian-nj/site/back/internal/validator"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu        sync.Mutex
		clientMap = make(map[string]*client)
	)
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clientMap {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clientMap, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			hIp, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			mu.Lock()
			if _, found := clientMap[hIp]; !found {
				clientMap[hIp] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}
			cl := clientMap[hIp]
			cl.lastSeen = time.Now()
			if !cl.limiter.Allow() {
				mu.Unlock()
				app.toManyRequestsResponse(w)
				return
			}
			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]
		v := validator.New()
		if data.ValidateTokenPlainText(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)

			return
		}
		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.invalidAuthenticationTokenResponse(w, r)

			} else {
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
