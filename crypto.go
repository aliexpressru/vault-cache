package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"io"
)

func getCacheKey(path, token string) string {
	var buf bytes.Buffer
	buf.WriteString(path)
	buf.WriteString(token)
	hash := sha256.Sum256(buf.Bytes())
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func encryptSecret(secret *Secret, keyStr string) ([]byte, error) {
	// Encode the Secret struct to a byte slice
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(secret)
	if err != nil {
		return nil, err
	}

	// Encrypt the byte slice
	encryptedData, err := encrypt(buf.Bytes(), keyStr)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}

func encrypt(data []byte, keyStr string) ([]byte, error) {
	key, err := generateAESKeyFromString(keyStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

func decryptSecret(encryptedData []byte, keyStr string) (*Secret, error) {
	// Decrypt the byte slice
	decryptedData, err := decrypt(encryptedData, keyStr)
	if err != nil {
		return nil, err
	}

	// Decode the byte slice into a Secret struct
	buf := bytes.NewBuffer(decryptedData)
	dec := gob.NewDecoder(buf)
	var secret Secret
	err = dec.Decode(&secret)
	if err != nil {
		return nil, err
	}

	return &secret, nil
}

func decrypt(data []byte, keyStr string) ([]byte, error) {
	key, err := generateAESKeyFromString(keyStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, errors.New("data is too short to be decrypted")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCTR(block, iv)
	decrypted := make([]byte, len(data))
	stream.XORKeyStream(decrypted, data)
	return decrypted, nil
}

func generateAESKeyFromString(key string) ([]byte, error) {
	if len(key) < aes.BlockSize {
		return nil, errors.New("key must be at least 16 bytes (128 bits)")
	}

	// Use the first 16 bytes of the key as the AES key
	aesKey := []byte(key)[:aes.BlockSize]

	// Generate a fixed IV (Initialization Vector)
	// We need fixed IV to achieve AES key determinism for one particular token
	iv := make([]byte, aes.BlockSize)
	copy(iv, "fixediv")

	// Create a new AES encryption block
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// Create a new CTR (Counter) mode of operation with the fixed IV
	cipher.NewCTR(block, iv)

	// Return the AES key and IV as a combined key
	return append(aesKey, iv...), nil
}
