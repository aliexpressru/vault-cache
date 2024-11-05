# Vault Proxy

This is a simple Vault proxy implementation in Go that provides an encrypted cache for Vault secrets.

## Features

- Transparently forwards all requests (except for `GET /v1/my-app-catalog/app/my-secret`) to the upstream Vault server.
- Caches Vault secrets in an in-memory cache, encrypted with the provided token.
- Automatically renews the cache when the TTL expires.
- Handles token-based authentication with the upstream Vault server.

## Run

To run the Vault proxy server, use the following command:
```shell
./vault-proxy [-upstream <upstream_url>] [-ttl <cache_ttl>] [-port <listen_port>]
```
where 
```shell

`-upstream`: The URL of the upstream Vault server. Default: `http://localhost:8200`.
`-ttl`: The TTL (Time-to-Live) for the cache. Default: `5m` (5 minutes).
`-port`: The port on which the Vault proxy server should listen. Default: `8201`.
```
For example, to run the Vault proxy server with an upstream Vault server at `https://vault.acme.com`, a cache TTL of 1 hour, and listening on port 8200, use the following command:
```shell
./vault-proxy -upstream https://vault.acme.com -ttl 1h -port 8200
```

## Usage

Send requests to the proxy server at `http://localhost:8200/v1/<path>`, providing the `X-Vault-Token` header with your Vault token.

Example:
```shell
curl -H "X-Vault-Token: your_token_here" http://localhost:8200/v1/my-app-catalog/app/my-secret
```


The proxy will fetch the secret from the upstream Vault server, encrypt it with the provided token, and store it in the cache. Subsequent requests for the same secret will be served from the cache.

## Testing

To run the tests for the `encrypt` and `decrypt` functions, use the following command:

```shell
go test -v
```
NB: do not forget to set up your tokens in `var` sequence of tests 


The tests ensure that the encryption and decryption functions work correctly and that the decrypted data matches the original input.

## Limitations

- The cache is stored in-memory and will be lost when the proxy server is restarted.
- The proxy does not handle token renewal or expiration. If the token used for encryption/decryption expires, the cached secrets will become inaccessible.
- The proxy does not provide any authentication or authorization mechanisms for clients. It relies on the Vault token provided in the request headers.


