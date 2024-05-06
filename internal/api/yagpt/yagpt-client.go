package yagpt

import (
	"context"

	grpc_client "github.com/Morwran/yagpt/internal/grpc-client"

	grpcClient "github.com/Morwran/yagpt/internal/3d-party/H-BF/corlib/client/grpc"
	"github.com/pkg/errors"
	ya "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/foundation_models/v1/text_generation"
	"google.golang.org/grpc"
)

type (
	// YaGPTClient
	YaGPTClient struct {
		ya.TextGenerationServiceClient
	}

	// ClosableClient
	YaGPTClosableClient struct {
		ya.TextGenerationServiceClient
		grpcClient.Closable
	}
)

// NewYaGPTClient constructs 'yagpt' API Client
func NewYaGPTClient(c grpc.ClientConnInterface) YaGPTClient {
	return YaGPTClient{
		TextGenerationServiceClient: ya.NewTextGenerationServiceClient(
			grpcClient.WithErrorWrapper(c, "yagpt"),
		),
	}
}

// NewYaGPTClosableClient constructs closable 'yagpt' API Client
func NewYaGPTClosableClient(ctx context.Context, p grpc_client.ConnProvider) (YaGPTClosableClient, error) {
	const api = "yagpt/new-closable-client"

	c, err := p.New(ctx)
	if err != nil {
		return YaGPTClosableClient{}, errors.WithMessage(err, api)
	}
	closable := grpcClient.MakeCloseable(
		grpcClient.WithErrorWrapper(c, "yagpt"),
	)
	return YaGPTClosableClient{
		TextGenerationServiceClient: ya.NewTextGenerationServiceClient(closable),
		Closable:                    closable,
	}, nil
}
