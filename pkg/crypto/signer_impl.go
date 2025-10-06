package crypto

import (
	stdlibcrypto "crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
)

// RSASigner wraps an RSA private key and signs payloads using PKCS#1 v1.5.
type RSASigner struct {
	key *rsa.PrivateKey
}

// NewRSASigner constructs an RSASigner from a private key.
func NewRSASigner(key *rsa.PrivateKey) *RSASigner {
	return &RSASigner{key: key}
}

// Sign signs the payload using SHA-256 and PKCS#1 v1.5.
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	if s == nil || s.key == nil {
		return nil, errors.New("rsa signer not initialised")
	}

	hash := sha256.Sum256(dataToBeSigned)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.key, stdlibcrypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("rsa sign: %w", err)
	}

	return signature, nil
}

// ECDSASigner signs using ECDSA and returns ASN.1 DER encoded signatures.
type ECDSASigner struct {
	key *ecdsa.PrivateKey
}

// NewECDSASigner constructs an ECDSASigner from a private key.
func NewECDSASigner(key *ecdsa.PrivateKey) *ECDSASigner {
	return &ECDSASigner{key: key}
}

// Sign signs the payload using SHA-256 and returns ASN.1 encoded signature.
func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	if s == nil || s.key == nil {
		return nil, errors.New("ecdsa signer not initialised")
	}

	hash := sha256.Sum256(dataToBeSigned)
	signature, err := ecdsa.SignASN1(rand.Reader, s.key, hash[:])
	if err != nil {
		return nil, fmt.Errorf("ecdsa sign: %w", err)
	}

	return signature, nil
}
