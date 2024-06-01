package yagpt

import (
	"context"
	"fmt"

	"github.com/Morwran/yagpt/internal"
	"github.com/pkg/errors"
	v1 "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/foundation_models/v1"
	ya "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/foundation_models/v1/text_generation"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const YaModelLite = "yandexgpt-lite"

type (
	YaGPTFace interface {
		CompletionWithCtx(ctx context.Context, iamTok string, m []Message) (*CompletionResponse, error)
		Completion(iamTok string, m []Message) (*CompletionResponse, error)
	}

	YaGPT struct {
		yaCli    *internal.YaGptClient
		folderId string
	}
)

func newYagpt(ctx context.Context, xfolderId string) (YaGPTFace, error) {
	cl, err := internal.NewYaGPTClient(ctx)
	if err != nil {
		return nil, err
	}
	return &YaGPT{yaCli: cl, folderId: xfolderId}, nil
}

func NewYagptWithCtx(ctx context.Context, xfolderId string) (YaGPTFace, error) {
	return newYagpt(ctx, xfolderId)
}

func NewYagpt(xfolderId string) (YaGPTFace, error) {
	ctx := context.Background()
	return newYagpt(ctx, xfolderId)
}

func (y *YaGPT) CompletionWithCtx(ctx context.Context, iamTok string, m []Message) (*CompletionResponse, error) {
	return y.completion(ctx, iamTok, m)
}

func (y *YaGPT) Completion(iamTok string, m []Message) (*CompletionResponse, error) {
	ctx := context.Background()
	return y.completion(ctx, iamTok, m)
}

func (y *YaGPT) completion(ctx context.Context, iamTok string, m []Message) (*CompletionResponse, error) {
	md := metadata.Pairs(
		"Authorization", "Bearer "+iamTok,
		"x-folder-id", y.folderId)

	ctx = metadata.NewOutgoingContext(ctx, md)
	var msgs []*v1.Message
	for _, msg := range m {
		msgs = append(msgs, msg.convertTo())
	}
	complCli, err := y.yaCli.Completion(ctx, &ya.CompletionRequest{
		ModelUri: fmt.Sprintf("gpt://%s/%s", y.folderId, YaModelLite),
		Messages: msgs,
		CompletionOptions: &v1.CompletionOptions{
			Temperature: &wrapperspb.DoubleValue{Value: 0.6},
			MaxTokens:   &wrapperspb.Int64Value{Value: 2000},
			Stream:      false,
		},
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed completion")
	}
	complResp, err := complCli.Recv()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to rcv response from yagpt")
	}

	var resp CompletionResponse
	return resp.convertFrom(complResp), nil
}
