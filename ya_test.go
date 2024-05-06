package yagpt

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type yaCliTestSuite struct {
	suite.Suite
	ctx context.Context
	iam IamFace
	ya  YaGPTFace
}

func (sui *yaCliTestSuite) SetupSuite() {
	sui.ctx = context.Background()
	oauthTok := os.Getenv("YA_OAUTH_TOK")
	sui.Require().NotEmpty(oauthTok)
	iam, err := NewYaIam(oauthTok)
	sui.Require().NoError(err)
	sui.iam = iam
	xfolderId := os.Getenv("YA_FOLDER_ID")
	sui.Require().NotEmpty(xfolderId)
	ya, err := NewYagptWithCtx(sui.ctx, xfolderId)
	sui.Require().NoError(err)
	sui.ya = ya
}

func (sui *yaCliTestSuite) Test_YaIam() {
	ctx := context.Background()

	resp, err := sui.iam.CreateWithCtx(ctx)
	sui.Require().NoError(err)
	sui.Require().NotNil(resp)
	sui.Require().NotEmpty(resp.IamToken)

}

func (sui *yaCliTestSuite) Test_YaCompl() {
	respIam, err := sui.iam.CreateWithCtx(sui.ctx)
	sui.Require().NoError(err)
	sui.Require().NotNil(respIam)
	sui.Require().NotEmpty(respIam.IamToken)

	var m []Message
	m = append(m, Message{Role: "user", Content: "hi"})
	resp, err := sui.ya.CompletionWithCtx(sui.ctx, respIam.IamToken, m)
	sui.Require().NoError(err)
	sui.Require().NotNil(resp)
	sui.Require().NotNil(resp.Alternatives)
	sui.Require().NotEmpty(resp.Alternatives[0].Message.Content)
}

func Test_YaCli(t *testing.T) {
	suite.Run(t, new(yaCliTestSuite))
}
