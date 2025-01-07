package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func publicKeyFromFile(keyFile string) (*rsa.PublicKey, error) {
	publicPem, _ := os.ReadFile(keyFile)

	block, _ := pem.Decode(publicPem)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to decode PEM block containing PKCS1 public key")
	}
	return publicKey, nil
}

func privateKeyFromFile(keyFile string) (*rsa.PrivateKey, error) {
	privatePem, _ := os.ReadFile(keyFile)

	block, _ := pem.Decode(privatePem)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to decode PEM block containing PKCS1 private key")
	}
	return privateKey, nil
}
