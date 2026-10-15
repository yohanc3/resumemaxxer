package llm 

import (
	"context"
)

type ResumeElements struct {
	Experiences []string
	ExperiencesBulletpoints map[string][]string
	Projects []string
	ProjectsBulletpoints map[string][]string
	skills map[string][]string
}

type LLM interface {
	CreateResumeElements(context.Context) (ResumeElements, error)
}
