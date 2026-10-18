package llm

import (
	"context"
)

type LLM interface {
	// Modifies user's resume experience/project descriptions and skills according to a job description
	//
	// It returns a [ResumeElements] object, which holds the new resume's information 
	CreateResumeElements(contex context.Context, resumeElements ResumeElements, jobDescription string, prompt string) (ResumeElements, error)

	// Parses a resume in string format into a structured format ([ResumeElements])
	ParseResume(context context.Context, resume string, prompt string) (ResumeElements, error)
}

type ResumeElements struct {
	Header     Header
	Eduaction  Education
	Experience []Experience
	Projects   []Projects
	Skills     map[string][]string
}

type Header struct {
	Name         string
	Email        string
	GithubLink   string
	LinkedinLink string
	PhoneNumber  string
}

type Education struct {
	School               string
	SchoolLocation       string // City, State if in the US, or City, Country if not
	DegreeType           string // Bachelors, Masters, PhD, etc.
	Degree               string // Computer Science, Philosophy, etc.
	ExpectedGraduation   string
	RelevantCoursework   []string
	LeadershipActivities []string
}

type Experience struct {
	CompanyName  string
	Title        string
	StartDate    string // [Month] [Year]
	EndDate      string // [Month] [Year]
	Location     string
	Bulletpoints []string
}

type Bulletpoint struct {
	Content string   // "Accomplished [X] as measured by [Y], by doing [Z]"
	Bolded  []string // ["Go/gRPC", "AWS Bedrock"]
}

type Projects struct {
	Projects             []string
	ProjectsBulletpoints map[string][]string
}
