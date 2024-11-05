package main

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

func writeSecret(w http.ResponseWriter, secret *Secret) {
	//write headers
	for key, values := range secret.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	//write response code
	w.WriteHeader(secret.StatusCode)
	//write body
	w.Write(secret.Body)
}

func getFromUpstream(r *http.Request, bodyBytes []byte) (*Secret, error) {
	path := r.URL.Path

	log.Debug().
		Str("path", path).
		Msg("Fetching secret from upstream")

	upstreamUrl := upstream + path
	// Create a new reader with the buffer content
	bodyReader := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(r.Method, upstreamUrl, bodyReader)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", path).
			Msg("Failed to create request to upstream")
		return nil, err
	}
	// Copy the headers from the original request
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	///work with result
	headers := map[string][]string{}
	for key, values := range resp.Header {
		headers[key] = values
	}

	statusCode := resp.StatusCode

	body, _ := ioutil.ReadAll(resp.Body)

	return &Secret{
		Headers:    headers,
		Body:       body,
		StatusCode: statusCode,
	}, nil

}
