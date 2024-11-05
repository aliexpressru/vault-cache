package main

import (
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	// Implement your health check logic here
	// Return a 200 HTTP status code if the application is healthy
	w.WriteHeader(http.StatusOK)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Increment the request counter
	vaultProxyRequests.Inc()

	path := r.URL.Path
	// Read the request body into a buffer
	// it could be closed before we want upstream it
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", path).
			Msg("Error reading body")
	}
	token, err := extractToken(r, bodyBytes)
	if err != nil {
		log.Error().Err(err).Msg("Error extraction token")
	} else {
		secret, err := getSecret(r, token)
		if err == nil {
			log.Info().
				Str("path", path).
				Msg("Served from cache")
			writeSecret(w, secret)
			return
		} else {
			log.Error().
				Err(err).
				Str("path", path).
				Msg("Cache error")
		}
	}

	secretEntry, err := getFromUpstream(r, bodyBytes)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", path).
			Msg("Fetch upstream")
	}

	cacheSecret(secretEntry, token, path)

	writeSecret(w, secretEntry)
}
