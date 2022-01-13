package tests

import (
	"github.com/rafaelcalleja/go-kit/internal/common/adapters"
	"github.com/rafaelcalleja/go-kit/internal/common/client"
	"github.com/stretchr/testify/require"
	"testing"
)

func NewStoreGrpcClient(t *testing.T, addr string) adapters.StoreGrpc {
	ok := WaitForPort(addr)
	require.True(t, ok, "Store Grpc timed out")

	storeClient, _, err := client.NewStoreClient(
		client.StoreClientWithGrpcAddress(addr),
		client.StoreClientWithNoTLS(true),
	)

	if err != nil {
		panic(err)
	}

	require.NoError(t, err)

	return adapters.NewStoreGrpc(storeClient)
}
