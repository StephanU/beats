// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"github.com/elastic/elastic-agent-libs/transport/tlscommon"
)

// NewProviderFromCertificate returns a TokenProvider that uses certificate-based
// authentication.
func NewProviderFromCertificate(resource, applicationID, tenantID string, conf tlscommon.CertificateConfig) (sptp TokenProvider, err error) {
	cert, privKey, err := loadConfigCerts(conf)
	if err != nil {
		return nil, fmt.Errorf("failed loading certificates: %w", err)
	}

	cred, err := azidentity.NewClientCertificateCredential(tenantID, applicationID, []*x509.Certificate{cert}, privKey, nil)
	if err != nil {
		return nil, err
	}

	return (*credentialTokenProvider)(cred), nil
}

func loadConfigCerts(cfg tlscommon.CertificateConfig) (cert *x509.Certificate, key *rsa.PrivateKey, err error) {
	tlsCert, err := tlscommon.LoadCertificate(&cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading X509 certificate from '%s': %w", cfg.Certificate, err)
	}
	if tlsCert == nil || len(tlsCert.Certificate) == 0 {
		return nil, nil, fmt.Errorf("no certificates loaded from '%s'", cfg.Certificate)
	}
	cert, err = x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing X509 certificate from '%s': %w", cfg.Certificate, err)
	}
	if tlsCert.PrivateKey == nil {
		return nil, nil, fmt.Errorf("failed loading private key from '%s'", cfg.Key)
	}
	key, ok := tlsCert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("private key at '%s' is not an RSA private key", cfg.Key)
	}
	return cert, key, nil
}
