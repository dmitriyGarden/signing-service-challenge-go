package crypto

import (
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
)

// DefaultKeyGenerator marshals generated key pairs into PEM encoded material.
type DefaultKeyGenerator struct {
	rsaGenerator RSAGenerator
	rsaMarshaler RSAMarshaler
	eccGenerator ECCGenerator
	eccMarshaler ECCMarshaler
}

var _ devices.KeyGenerator = (*DefaultKeyGenerator)(nil)

// NewDefaultKeyGenerator instantiates a generator for RSA and ECDSA keys.
func NewDefaultKeyGenerator() *DefaultKeyGenerator {
	return &DefaultKeyGenerator{
		rsaGenerator: RSAGenerator{},
		rsaMarshaler: NewRSAMarshaler(),
		eccGenerator: ECCGenerator{},
		eccMarshaler: NewECCMarshaler(),
	}
}

// Generate produces PEM encoded key material for the requested algorithm.
func (g *DefaultKeyGenerator) Generate(algorithm domain.Algorithm) (domain.KeyMaterial, error) {
	switch algorithm {
	case domain.AlgorithmRSA:
		return g.generateRSA()
	case domain.AlgorithmECDSA:
		return g.generateECDSA()
	default:
		return domain.KeyMaterial{}, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func (g *DefaultKeyGenerator) generateRSA() (domain.KeyMaterial, error) {
	pair, err := g.rsaGenerator.Generate()
	if err != nil {
		return domain.KeyMaterial{}, fmt.Errorf("generate rsa key pair: %w", err)
	}

	publicBytes, privateBytes, err := g.rsaMarshaler.Marshal(*pair)
	if err != nil {
		return domain.KeyMaterial{}, fmt.Errorf("marshal rsa key pair: %w", err)
	}

	return domain.KeyMaterial{
		// Copy slices so callers cannot mutate the generator's buffers.
		Public:  append([]byte(nil), publicBytes...),
		Private: append([]byte(nil), privateBytes...),
	}, nil
}

func (g *DefaultKeyGenerator) generateECDSA() (domain.KeyMaterial, error) {
	pair, err := g.eccGenerator.Generate()
	if err != nil {
		return domain.KeyMaterial{}, fmt.Errorf("generate ecdsa key pair: %w", err)
	}

	publicBytes, privateBytes, err := g.eccMarshaler.Encode(*pair)
	if err != nil {
		return domain.KeyMaterial{}, fmt.Errorf("marshal ecdsa key pair: %w", err)
	}

	return domain.KeyMaterial{
		// Copy slices so callers cannot mutate the generator's buffers.
		Public:  append([]byte(nil), publicBytes...),
		Private: append([]byte(nil), privateBytes...),
	}, nil
}
