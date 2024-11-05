package main

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Secret struct {
	Headers    map[string][]string
	Body       []byte
	StatusCode int
}

func extractToken(r *http.Request, bodyBytes []byte) (string, error) {
	// GET -> extract X-Vault-Token
	if r.Method == http.MethodGet {
		token := r.Header.Get("X-Vault-Token")
		return token, nil
	}
	// POST -> extract JWT
	if r.Method == http.MethodPost {
		var jsonBody map[string]interface{}
		err := json.Unmarshal(bodyBytes, &jsonBody)
		if err != nil {
			log.Error().
				Err(err).
				Msg("extract token error")
			return "", err
		}
		jwtToken := jsonBody["jwt"]
		return jwtToken.(string), nil
	}
	return "", fmt.Errorf("no token found")
}

func cacheSecret(secretEntry *Secret, token string, path string) {
	encryptedSecret, err := encryptSecret(secretEntry, token)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", path).
			Msg("Failed to encrypt secret")
		return
	}
	cacheKey := getCacheKey(path, token)
	c.Set(cacheKey, encryptedSecret, cache.DefaultExpiration)
	log.Debug().
		Str("key", cacheKey).
		Str("path", path).
		Msg("Cached")
}

func getSecret(r *http.Request, token string) (*Secret, error) {
	if token != "" {
		path := r.URL.Path
		cacheKey := getCacheKey(path, token)
		encryptedSecret, ok := c.Get(cacheKey)
		if ok {
			secret, err := decryptSecret(encryptedSecret.([]byte), token)
			if err != nil {
				log.Error().
					Err(err).
					Str("path", path).
					Msg("Failed to decrypt secret")
				return nil, err
			}
			cacheHits.Inc()
			return secret, nil
		} else {
			return nil, fmt.Errorf("not in cache")
		}
	}
	return nil, fmt.Errorf("empty token")
}
