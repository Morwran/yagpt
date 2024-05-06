package internal

import (
	"context"

	yaAPI "github.com/Morwran/yagpt/internal/api/yagpt"
	grpc_client "github.com/Morwran/yagpt/internal/grpc-client"
	"github.com/Morwran/yagpt/internal/tls"
	"github.com/pkg/errors"
)

// YaGptClient is an alias to 'yaAPI.ClosableClient'
type YaGptClient = yaAPI.YaGPTClosableClient

// NewYaGPTClient makes 'yagpt' API client
func NewYaGPTClient(ctx context.Context) (*YaGptClient, error) {
	clientCreds, err := tls.ClientTransportCredentials(true, "", "", "")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create client creds")
	}
	bld := grpc_client.FromAddress(yagptUrl).
		WithDialDuration(defDialDuration).
		WithCreds(clientCreds)

	c, err := yaAPI.NewYaGPTClosableClient(ctx, bld)
	if err != nil {
		return nil, err
	}
	return &c, err
}
