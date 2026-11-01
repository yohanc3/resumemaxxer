package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	llm "github.com/yohanc3/resumemaxxer/internal/llm"
)

type DeepseekClient struct {
	client                        openai.Client
	llmResumeElementsFormatSchema map[string]any
}

func NewDeepseekClient(baseURL string, apiKey string) (*DeepseekClient, error) {
	resumeElementsJSONSchema, err := llm.GenerateSchema[llm.ResumeElements]()
	if err != nil {
		return nil, fmt.Errorf("error when generating resume elements schema at deepseek client. %w", err)
	}

	return &DeepseekClient{
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
			option.WithBaseURL(baseURL),
		),
		llmResumeElementsFormatSchema: resumeElementsJSONSchema,
	}, nil
}

// Todo
func (d *DeepseekClient) CreateResumeElements(ctx context.Context, resumeElements *llm.ResumeElements, jobDescription string, prompt string)  (*llm.ResumeElements, error) {
	return &llm.ResumeElements{}, nil
}

func (d *DeepseekClient) GetResumeElements(ctx context.Context, resume string, prompt string) (*llm.ResumeElements, error) {

	resp, err := d.client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
		Model: "deepseek-flash",
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name: "extracted_resume_elements",
					Description: openai.String(`
					Elements from a resume such as:
						- Header (first name, last name, email, phone number, linkedin link, github link),
						- Education (school, school location, degree type, degree,
									 expected graduation in [month], [year] format, relevant coursework,
									 and leadership activities)
						- Experiences (company name, title, start date, end date, location, bulletpoints) 
						- Projects (project name, role, bulletpoints)
						- Bulletpoint - used in experiences and projects:  (TextSegments) 
						- TextSegment: - used in bulletpoints to mark bulletpoints text bolded (text, bold?)
					`),
					Schema: d.llmResumeElementsFormatSchema,
				},
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error when calling AI to parse resume elements: %w", err)
	}

	var resumeElements llm.ResumeElements

	if err := json.Unmarshal([]byte(resp.OutputText()), &resumeElements); err != nil {
		return nil, fmt.Errorf("error when unmarshaling AI parsed resume elements response: %w", err)
	}

	return &resumeElements, nil
}
