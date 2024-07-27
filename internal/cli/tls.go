package cli

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	"github.com/sebasttiano/Owl/internal/logger"
	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials(caCert string) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile(caCert)
	if err != nil {
		logger.Log.Error("failed to load certificate of the CA")
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		logger.Log.Error("failed to add server CA's certificate")
		return nil, errors.New("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}
