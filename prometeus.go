package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	vaultProxyRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vault_cache_requests_total",
		Help: "Total number of requests handled by the Vault proxy",
	})

	cacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vault_cache_cache_hits",
		Help: "Cache hits encountered by the Vault proxy",
	})
)
