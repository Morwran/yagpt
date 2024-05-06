package tls

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
)

func ClientTLSConfig(insecureSkipVerify bool, cacertFile, clientCertFile, clientKeyFile string) (*tls.Config, error) {
	var tlsConf tls.Config

	if clientCertFile != "" {
		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return nil, errors.WithMessage(err, "could not load client key pair")
		}
		tlsConf.Certificates = []tls.Certificate{certificate}
	}

	if insecureSkipVerify {
		tlsConf.InsecureSkipVerify = true
	} else if cacertFile != "" {
		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := os.ReadFile(cacertFile)
		if err != nil {
			return nil, errors.WithMessage(err, "could not read ca certificate")
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("failed to append ca certs")
		}

		tlsConf.RootCAs = certPool
	}

	return &tlsConf, nil
}

func ClientTransportCredentials(insecureSkipVerify bool, cacertFile, clientCertFile, clientKeyFile string) (credentials.TransportCredentials, error) {
	tlsConf, err := ClientTLSConfig(insecureSkipVerify, cacertFile, clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(tlsConf), nil
}
