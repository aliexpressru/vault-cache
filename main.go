package main

import (
	"flag"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

var (
	upstream string
	ttl      time.Duration
	port     int
	c        *cache.Cache
)

func main() {
	// Parse command-line arguments
	flag.StringVar(&upstream, "upstream", "http://localhost:8200", "URL of the upstream Vault server")
	flag.DurationVar(&ttl, "ttl", 5*time.Minute, "TTL for the cache")
	flag.IntVar(&port, "port", 8201, "Port to listen on")
	flag.Parse()

	if port > 65535 || port < 1 {
		log.Fatal().
			Int("port", port).
			Msg("Wrong port")
	}
	log.Info().
		Str("upstream", upstream).
		Float64("ttl", ttl.Seconds()).
		Int("listen port", port).
		Msg("Received args")

	c = cache.New(ttl, ttl*2)
	log.Info().Msg("Cache initialized")

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	// Set the minimum log level to debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/healthz", healthzHandler)
	router.HandleFunc("/v1/{path:.+}", proxyHandler).Methods(http.MethodGet, http.MethodPost)

	log.Info().Msg("Starting Vault proxy server")
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
