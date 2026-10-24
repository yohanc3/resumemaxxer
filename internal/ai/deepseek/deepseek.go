package ai

import (
	"context"

	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/yohanc3/resumemaxxer/internal/ai/llm"
)

type DeepseekClient struct {
	client openai.Client 
}

func NewDeepseekClient(baseURL string, apiKey string) *DeepseekClient {
	return &DeepseekClient{client: openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)}
}


// Todo
func (d *DeepseekClient) CreateResumeElements(ctx context.Context, resume string, prompt string) (*llm.ResumeElements, error){
	return &llm.ResumeElements{}, nil
}


func (d *DeepseekClient) ParseResume(ctx context.Context, resume string, prompt string) (*llm.ResumeElements, error){

	_, _ = d.client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
		Model: "deepseek-flash",  
	})

	return nil, nil

}
