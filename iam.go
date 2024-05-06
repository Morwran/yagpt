package yagpt

import (
	"context"
	"time"

	"github.com/Morwran/yagpt/internal"
	"github.com/pkg/errors"
	iam "github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
)

type (
	IamTokenResponse struct {

		// IAM token for the specified identity.
		//
		// You should pass the token in the `Authorization` header for any further API requests.
		// For example, `Authorization: Bearer [iam_token]`.
		IamToken string
		// IAM token expiration time.
		ExpiresAt time.Time
	}
	IamFace interface {
		Create() (*IamTokenResponse, error)
		CreateWithCtx(ctx context.Context) (*IamTokenResponse, error)
		Close() error
	}
	iamImpl struct {
		iamCli *internal.YaIamClient
		tok    string
	}
)

func newIam(ctx context.Context, oauthTok string) (IamFace, error) {
	cl, err := internal.NewYaIamClient(ctx)
	if err != nil {
		return nil, err
	}
	return &iamImpl{iamCli: cl, tok: oauthTok}, nil
}

func NewYaIamWithCtx(ctx context.Context, oauthTok string) (IamFace, error) {
	return newIam(ctx, oauthTok)
}

func NewYaIam(oauthTok string) (IamFace, error) {
	ctx := context.Background()
	return newIam(ctx, oauthTok)
}

func (i iamImpl) create(ctx context.Context) (*IamTokenResponse, error) {
	resp, err := i.iamCli.Create(ctx, &iam.CreateIamTokenRequest{
		Identity: &iam.CreateIamTokenRequest_YandexPassportOauthToken{
			YandexPassportOauthToken: i.tok,
		},
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to request iam tok")
	}
	t := resp.GetExpiresAt()
	if t == nil {
		return nil, errors.WithMessage(err, "failed to get expires time")
	}
	return &IamTokenResponse{
		IamToken:  resp.GetIamToken(),
		ExpiresAt: t.AsTime()}, nil
}

func (i iamImpl) CreateWithCtx(ctx context.Context) (*IamTokenResponse, error) {
	return i.create(ctx)
}

func (i iamImpl) Create() (*IamTokenResponse, error) {
	ctx := context.Background()
	return i.create(ctx)
}

func (i iamImpl) Close() error {
	return i.iamCli.CloseConn()
}
