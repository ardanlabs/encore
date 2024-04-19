// Package keystore implements the auth.KeyLookup interface. This implements
// an in-memory keystore for JWT support.
package keystore

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// key represents key information.
type key struct {
	privatePEM string
	publicPEM  string
}

// KeyStore represents an in memory store implementation of the
// KeyLookup interface for use with the auth package.
type KeyStore struct {
	store map[string]key
}

// New constructs an empty KeyStore ready for use.
func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]key),
	}
}

// LoadKey takes an id and the private PEM string.
func (ks *KeyStore) LoadKey(id string, pem string) error {
	privatePEM := string(pem)
	publicPEM, err := toPublicPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("converting private PEM to public: %w", err)
	}

	key := key{
		privatePEM: privatePEM,
		publicPEM:  publicPEM,
	}

	ks.store[id] = key

	return nil
}

// PrivateKey searches the key store for a given kid and returns the private key.
func (ks *KeyStore) PrivateKey(kid string) (string, error) {
	key, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}

	return key.privatePEM, nil
}

// PublicKey searches the key store for a given kid and returns the public key.
func (ks *KeyStore) PublicKey(kid string) (string, error) {
	key, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}

	return key.publicPEM, nil
}

func toPublicPEM(privatePEM string) (string, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return "", errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	var parsedKey any
	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
	}

	pk, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("key is not a valid RSA private key")
	}

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return "", fmt.Errorf("marshaling public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &publicBlock); err != nil {
		return "", fmt.Errorf("encoding to public PEM: %w", err)
	}

	return buf.String(), nil
}
