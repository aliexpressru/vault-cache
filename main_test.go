package main

import (
	"bytes"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	XVaultToken = ""
	jwt         = ""
	role        = "my-app"
)

func setup() {
	ttl = time.Duration(5 * time.Minute)
	upstream = "http://localhost:8200"
	port = 8201
	c = cache.New(ttl, ttl*2)
}

func TestGetProxy(t *testing.T) {
	setup()
	// Create a new request to simulate a GET request to a specific path
	c = cache.New(ttl, ttl*2)
	path := fmt.Sprintf("http://localhost:8200/v1/app/%s", "my-catalog/dev/my-app/value1")
	req := httptest.NewRequest("GET", path, nil)
	// TODO: set your token here, if you want test
	req.Header.Set("X-Vault-Token", XVaultToken)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the proxyHandler function directly with the simulated request and response recorder
	proxyHandler(rr, req)
	proxyHandler(rr, req)
	proxyHandler(rr, req)
	proxyHandler(rr, req)
	proxyHandler(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestPostProxy(t *testing.T) {
	setup()
	jsonPayload := []byte(fmt.Sprintf(`{"role":"%s","jwt":"%s"}`, role, jwt))
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the proxyHandler function directly with the simulated request and response recorder
	proxyHandler(rr, httptest.NewRequest("POST", "http://localhost:8200/v1/data/auth", bytes.NewBuffer(jsonPayload)))
	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	proxyHandler(rr, httptest.NewRequest("POST", "http://localhost:8200/v1/data/auth", bytes.NewBuffer(jsonPayload)))
	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	proxyHandler(rr, httptest.NewRequest("POST", "http://localhost:8200/v1/data/auth", bytes.NewBuffer(jsonPayload)))
	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	proxyHandler(rr, httptest.NewRequest("POST", "http://localhost:8200/v1/data/auth", bytes.NewBuffer(jsonPayload)))
	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	proxyHandler(rr, httptest.NewRequest("POST", "http://localhost:8200/v1/data/auth", bytes.NewBuffer(jsonPayload)))
	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestEncryptDecrypt(t *testing.T) {
	inputStrings := []string{"hello"}
	//inputStrings := []string{"hello", "world", "test", "encryption"}
	key := "mysecretkeydfdfdfdfdfdfdfdf"

	for _, input := range inputStrings {
		// Encrypt the input string
		log.Info().
			Bytes("Data", []byte(input)).
			Msg("Using data")
		encrypted, err := encrypt([]byte(input), key)
		if err != nil {
			t.Errorf("Failed to encrypt string: %s", err)
			continue
		}

		// Decrypt the encrypted string
		decrypted, err := decrypt(encrypted, key)
		if err != nil {
			t.Errorf("Failed to decrypt string: %s", err)
			continue
		}
		log.Info().
			Bytes("data", decrypted).
			Msg("Decrypted")
		// Check that the decrypted string matches the original input string
		if string(decrypted) != input {
			t.Errorf("Decrypted string does not match input string: %s != %s", decrypted, input)
		}
	}
}
