package gpt

import (
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"strings"
)

type TikTokenClient struct {
	tikToken *tiktoken.Tiktoken
}

func NewTikTokenClient(model string) (*TikTokenClient, error) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return nil, fmt.Errorf("encoding for model: %v", err)
	}
	return &TikTokenClient{
		tikToken: tkm,
	}, nil
}

// NumTokensFromMessages
// OpenAI Cookbook: https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
func (c *TikTokenClient) NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int, err error) {
	var tokensPerMessage, tokensPerName int
	switch model {
	case openai.GPT432K0613,
		openai.GPT432K0314,
		openai.GPT432K,
		openai.GPT40613,
		openai.GPT40314,
		openai.GPT4TurboPreview,
		openai.GPT4VisionPreview,
		openai.GPT4,
		openai.GPT3Dot5Turbo1106,
		openai.GPT3Dot5Turbo16K,
		openai.GPT3Dot5Turbo16K0613,
		openai.GPT3Dot5Turbo0613:
		tokensPerMessage = 3
		tokensPerName = 1
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
		tokensPerName = -1   // if there's a name, the role is omitted
	default:
		if strings.Contains(model, "gpt-3.5-turbo") {
			return c.NumTokensFromMessages(messages, "gpt-3.5-turbo-0613")
		} else if strings.Contains(model, "gpt-4") {
			return c.NumTokensFromMessages(messages, "gpt-4-0613")
		} else {
			err = fmt.Errorf("num_tokens_from_messages() is not implemented for model %s. See https://github.com/openai/openai-python/blob/main/chatml.md for information on how messages are converted to tokens", model)
			return
		}
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(c.tikToken.Encode(message.Content, nil, nil))
		numTokens += len(c.tikToken.Encode(message.Role, nil, nil))
		numTokens += len(c.tikToken.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>
	return numTokens, nil
}

func (c *TikTokenClient) Encode(content string) []int {
	return c.tikToken.Encode(content, nil, nil)
}

func (c *TikTokenClient) Decode(tokens []int) string {
	return c.tikToken.Decode(tokens)
}
