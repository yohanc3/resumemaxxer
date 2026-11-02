package llm

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

// Structured Outputs uses a subset of JSON schema
// These flags are necessary to comply with the subset
// Source: https://github.com/openai/openai-go under "Structured Outputs"
func GenerateSchema[T any]() (map[string]any, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)

	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	// Preserve schema values as raw JSON so integer constraints are not rounded
	// through float64 before the SDK serializes the request.
	var rawSchema map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawSchema); err != nil {
		return nil, err
	}
	result := make(map[string]any, len(rawSchema))
	for key, value := range rawSchema {
		result[key] = value
	}
	return result, nil
}

const RESUME_EXTRACTION_PROMPT = `
	Your job is to extract the elements of the given resume into its parts. The format for your response
	is in the JSON schema format, which contains the elements you must extract from the resume.

	The resume text given to you may be compact, but it is expected for the information to still be 
	readable enough to be identified.
`
