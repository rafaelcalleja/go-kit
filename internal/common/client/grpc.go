package client

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/rafaelcalleja/go-kit/internal/common/genproto/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

type StoreClientOptions struct {
	grpcAddr string
	noTLS    bool
}

func StoreClientWithGrpcAddress(grpcAddr string) func(*StoreClientOptions) error {
	return func(o *StoreClientOptions) error {
		o.grpcAddr = grpcAddr
		return nil
	}
}

func StoreClientWithNoTLS(noTLS bool) func(*StoreClientOptions) error {
	return func(o *StoreClientOptions) error {
		o.noTLS = noTLS
		return nil
	}
}

func NewStoreClient(options ...func(*StoreClientOptions) error) (client store.StoreServiceClient, close func() error, err error) {
	var storeClientOptions = new(StoreClientOptions)

	for _, option := range options {
		err := option(storeClientOptions)
		if err != nil {
			return nil, func() error { return nil }, err
		}
	}

	grpcAddr := storeClientOptions.grpcAddr

	if grpcAddr == "" {
		return nil, func() error { return nil }, errors.New("empty StoreClientOptions.grpcAddr")
	}

	opts, err := grpcDialOpts(grpcAddr, storeClientOptions.noTLS)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return store.NewStoreServiceClient(conn), conn.Close, nil
}

func WaitForStoreService(address string, timeout time.Duration) bool {
	return waitForPort(address, timeout)
}

func grpcDialOpts(grpcAddr string, noTLS bool) ([]grpc.DialOption, error) {
	if noTLS {
		return []grpc.DialOption{grpc.WithInsecure()}, nil
	}

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "cannot load root CA cert")
	}

	creds := credentials.NewTLS(&tls.Config{
		RootCAs:    systemRoots,
		MinVersion: tls.VersionTLS12,
	})

	return []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(newMetadataServerToken(grpcAddr)),
	}, nil
}
