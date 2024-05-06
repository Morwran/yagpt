[![Go Reference](https://pkg.go.dev/badge/github.com/Morwran/yagpt.svg)](https://pkg.go.dev/github.com/Morwran/yagpt)
[![Go Report Card](https://goreportcard.com/badge/github.com/Morwran/yagpt)](https://goreportcard.com/report/github.com/Morwran/yagpt)


## Go YandexGPT
#### This library provides simple Go clients for [Yandex Foundation Models](https://yandex.cloud/ru/docs/foundation-models/)

#### We support:
- YandexGPT Lite model
- Connect over GRPC
- IAM token generation

## Installation
    go get github.com/Morwran/yagpt

## Usage
### YaGPT example usage with generating AIM token:
```go
package main

import (
	"context"
	"fmt"
	"github.com/Morwran/yagpt"
)

func main() {
    ctx := context.Background()
    oauthTok := os.Getenv("YA_OAUTH_TOK")
    iam, err := yagpt.NewYaIam(oauthTok)
    if err != nil {
		fmt.Printf("filed to create connection for iam: %v\n", err)
		return
	}
    respIam, err := iam.CreateWithCtx(ctx)
    if err != nil {
		fmt.Printf("filed to generating iam token: %v\n", err)
		return
	}
    xfolderId := os.Getenv("YA_FOLDER_ID")
    ya, err := yagpt.NewYagptWithCtx(ctx, xfolderId)
    if err != nil {
		fmt.Printf("filed to create connection for yagpt: %v\n", err)
		return
	}
    var m []yagpt.Message
	m = append(m, yagpt.Message{Role: "user", Content: "hi"})
	resp, err := ya.CompletionWithCtx(ctx, respIam.IamToken, m)
    if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
    fmt.Println(resp.Alternatives[0].Message.Content)
}
```
### Getting an OAuth Yandex token and folder Id:

- [Authorize](https://passport.yandex.ru/auth) Yandex account
- Get [OAuth token](https://oauth.yandex.ru/)
- Get [folder id](https://yandex.cloud/ru/docs/resource-manager/operations/folder/get-id)