package yagpt

import (
	v1 "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/foundation_models/v1"
	ya "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/foundation_models/v1/text_generation"
)

type (
	Message struct {
		// A role used by the user to describe requests to the model.
		Role string
		// Message content.
		Content string
	}

	// Represents a generated completion alternative, including its content and generation status.
	Alternative struct {
		// A message containing the content of the alternative.
		Message Message
		// The generation status of the alternative
		Status v1.Alternative_AlternativeStatus
	}

	ContentUsage struct {
		// The number of tokens in the textual part of the model input.
		InputTextTokens int64
		// The total number of tokens in the generated completions.
		CompletionTokens int64
		// The total number of tokens, including all input tokens and all generated tokens.
		TotalTokens int64
	}

	CompletionResponse struct {
		// A list of generated completion alternatives.
		Alternatives []Alternative
		// A set of statistics describing the number of content tokens used by the completion model.
		Usage ContentUsage
		// The model version changes with each new releases.
		ModelVersion string
	}
)

func (m *Message) convertTo() *v1.Message {
	return &v1.Message{Role: m.Role, Content: &v1.Message_Text{Text: m.Content}}
}

func (m *Message) convertFrom(gMsg *v1.Message) *Message {
	if gMsg == nil {
		return m
	}
	m.Role = gMsg.GetRole()
	m.Content = gMsg.GetText()
	return m
}

func (r *CompletionResponse) convertFrom(resp *ya.CompletionResponse) *CompletionResponse {
	if resp == nil {
		return r
	}
	r.ModelVersion = resp.GetModelVersion()
	u := resp.GetUsage()
	if u != nil {
		r.Usage.CompletionTokens = u.GetCompletionTokens()
		r.Usage.InputTextTokens = u.GetInputTextTokens()
		r.Usage.TotalTokens = u.GetTotalTokens()
	}
	a := resp.GetAlternatives()
	var altTmp []Alternative
	for _, alt := range a {
		m := alt.GetMessage()
		var msg Message

		altTmp = append(altTmp, Alternative{
			Message: *msg.convertFrom(m),
			Status:  alt.GetStatus()})
	}
	r.Alternatives = append(r.Alternatives, altTmp...)
	return r
}
