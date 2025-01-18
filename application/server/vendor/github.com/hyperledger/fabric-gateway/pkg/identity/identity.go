// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package identity defines a client identity and signing implementation used to interact with a Fabric network.
//
// This package provides utilities to aid creation of client identities and accompanying signing implementations from
// various types of credentials.
package identity

import (
	"crypto/x509"
)

// Identity represents a client identity used to interact with a Fabric network.
type Identity interface {
	MspID() string       // ID of the Membership Service Provider to which this identity belongs.
	Credentials() []byte // Implementation-specific credentials.
}

// X509Identity represents a client identity backed by an X.509 certificate.
type X509Identity struct {
	mspID       string
	certificate []byte
}

// MspID returns the ID of the Membership Service Provider to which this identity belongs.
func (id *X509Identity) MspID() string {
	return id.mspID
}

// Credentials as an X.509 certificate in PEM encoded ASN.1 DER format.
func (id *X509Identity) Credentials() []byte {
	return id.certificate
}

// NewX509Identity creates a new Identity from an X.509 certificate.
func NewX509Identity(mspID string, certificate *x509.Certificate) (*X509Identity, error) {
	certificatePEM, err := CertificateToPEM(certificate)
	if err != nil {
		return nil, err
	}

	identity := &X509Identity{
		mspID:       mspID,
		certificate: certificatePEM,
	}
	return identity, nil
}
