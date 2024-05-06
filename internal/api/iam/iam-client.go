package iam

import (
	"context"

	grpc_client "github.com/Morwran/yagpt/internal/grpc-client"

	grpcClient "github.com/Morwran/yagpt/internal/3d-party/H-BF/corlib/client/grpc"
	"github.com/pkg/errors"
	iam "github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
	"google.golang.org/grpc"
)

type (
	// YaIamClient
	YaIamClient struct {
		iam.IamTokenServiceClient
	}

	// YaIamClosableClient
	YaIamClosableClient struct {
		iam.IamTokenServiceClient
		grpcClient.Closable
	}
)

// NewYaIamClient constructs 'iam' API Client
func NewYaIamClient(c grpc.ClientConnInterface) YaIamClient {
	return YaIamClient{
		IamTokenServiceClient: iam.NewIamTokenServiceClient(
			grpcClient.WithErrorWrapper(c, "iam"),
		),
	}
}

// NewYaIamClosableClient constructs closable 'iam' API Client
func NewYaIamClosableClient(ctx context.Context, p grpc_client.ConnProvider) (YaIamClosableClient, error) {
	const api = "iam/new-closable-client"

	c, err := p.New(ctx)
	if err != nil {
		return YaIamClosableClient{}, errors.WithMessage(err, api)
	}
	closable := grpcClient.MakeCloseable(
		grpcClient.WithErrorWrapper(c, "iam"),
	)
	return YaIamClosableClient{
		IamTokenServiceClient: iam.NewIamTokenServiceClient(closable),
		Closable:              closable,
	}, nil
}
