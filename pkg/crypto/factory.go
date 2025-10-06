package crypto

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
)

// SignerFactory resolves crypto signers based on the device algorithm.
type SignerFactory struct{}

var _ devices.SignerFactory = (*SignerFactory)(nil)

// NewSignerFactory instantiates a default factory.
func NewSignerFactory() *SignerFactory {
	return &SignerFactory{}
}

// SignerFor decodes key material and returns the matching signer.
func (f *SignerFactory) SignerFor(device domain.Device, material domain.KeyMaterial) (devices.Signer, error) {
	switch device.Algorithm {
	case domain.AlgorithmRSA:
		privateKey, err := parseRSAPrivateKey(material.Private)
		if err != nil {
			return nil, fmt.Errorf("decode rsa private key: %w", err)
		}
		return NewRSASigner(privateKey), nil
	case domain.AlgorithmECDSA:
		privateKey, err := parseECDSAPrivateKey(material.Private)
		if err != nil {
			return nil, fmt.Errorf("decode ecdsa private key: %w", err)
		}
		return NewECDSASigner(privateKey), nil
	default:
		return nil, domain.ErrInvalidAlgorithm
	}
}

// parseRSAPrivateKey extracts a PKCS#1 private key from a PEM encoded block.
func parseRSAPrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM block")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// parseECDSAPrivateKey extracts an EC private key from a PEM encoded block.
func parseECDSAPrivateKey(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM block")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
