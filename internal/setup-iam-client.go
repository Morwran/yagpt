package internal

import (
	"context"

	iamAPI "github.com/Morwran/yagpt/internal/api/iam"
	grpc_client "github.com/Morwran/yagpt/internal/grpc-client"
	"github.com/Morwran/yagpt/internal/tls"
	"github.com/pkg/errors"
)

// YaIamClient is an alias to 'iamAPI.ClosableClient'
type YaIamClient = iamAPI.YaIamClosableClient

// NewYaIamClient makes 'iam' API client
func NewYaIamClient(ctx context.Context) (*YaIamClient, error) {
	clientCreds, err := tls.ClientTransportCredentials(true, "", "", "")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create client creds")
	}
	bld := grpc_client.FromAddress(iamUri).
		WithDialDuration(defDialDuration).
		WithCreds(clientCreds)

	c, err := iamAPI.NewYaIamClosableClient(ctx, bld)
	if err != nil {
		return nil, err
	}
	return &c, err
}
