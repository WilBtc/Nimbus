// server/src/utils/encryption_utils.go

package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// EncryptData encrypts data using the provided RSA public key with OAEP padding and SHA-256 hashing.
func EncryptData(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}
	return ciphertext, nil
}

// DecryptData decrypts data using the provided RSA private key with OAEP padding and SHA-256 hashing.
func DecryptData(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}
	return plaintext, nil
}

// GenerateKeyPair generates a new RSA key pair with the specified bit size.
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("key generation failed: %w", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

// SavePrivateKey saves an RSA private key to a file in PEM format.
func SavePrivateKey(privateKey *rsa.PrivateKey, filename string) error {
	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer file.Close()

	err = pem.Encode(file, pemBlock)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}
	return nil
}

// SavePublicKey saves an RSA public key to a file in PEM format.
func SavePublicKey(publicKey *rsa.PublicKey, filename string) error {
	keyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %w", err)
	}
	defer file.Close()

	err = pem.Encode(file, pemBlock)
	if err != nil {
		return fmt.Errorf("failed to encode public key: %w", err)
	}
	return nil
}

// SignData creates a digital signature for the given data using the provided RSA private key and SHA-512 hashing.
func SignData(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	hashed := sha512.Sum512(data)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA512, hashed[:], nil)
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}
	return signature, nil
}

// VerifySignature verifies a digital signature using the provided RSA public key and SHA-512 hashing.
func VerifySignature(publicKey *rsa.PublicKey, data, signature []byte) error {
	hashed := sha512.Sum512(data)
	err := rsa.VerifyPSS(publicKey, crypto.SHA512, hashed[:], signature, nil)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}
	return nil
}
