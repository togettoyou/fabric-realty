// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// CertificateToPEM converts an X.509 certificate to PEM encoded ASN.1 DER data.
func CertificateToPEM(certificate *x509.Certificate) ([]byte, error) {
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	}
	return pemEncode(block)
}

// CertificateFromPEM creates an X.509 certificate from PEM encoded data.
func CertificateFromPEM(certificatePEM []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certificatePEM)
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}

	return x509.ParseCertificate(block.Bytes)
}

// PrivateKeyFromPEM creates a private key from PEM encoded data.
func PrivateKeyFromPEM(privateKeyPEM []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse private key PEM")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// PrivateKeyToPEM converts a private key to PEM encoded PKCS #8 data.
func PrivateKeyToPEM(privateKey crypto.PrivateKey) ([]byte, error) {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pemEncode(block)
}

func pemEncode(block *pem.Block) ([]byte, error) {
	var buffer bytes.Buffer
	if err := pem.Encode(&buffer, block); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
