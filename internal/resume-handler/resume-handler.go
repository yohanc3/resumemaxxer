package resumehandler

import (
	"context"
	"database/sql"

	"github.com/yohanc3/resumemaxxer/internal/llm"
	objectstorage "github.com/yohanc3/resumemaxxer/internal/storage/objectStorage"
)

type QueueJob struct {
	ID                    string
	JobPostingID          string
	JobPostingDescription string
	UserID                string
	Retries               int
	Status                string
}

type ResumeHandler struct {
	db            *sql.DB
	objectStorage objectstorage.ObjectStorage
	resumeLLM     llm.ResumeLLM
}

func NewResumeHandler(db *sql.DB, objectStorage objectstorage.ObjectStorage, llm llm.ResumeLLM) ResumeHandler {
	return ResumeHandler{db: db, objectStorage: objectStorage, resumeLLM: llm}
}

func (r *ResumeHandler) CreateResume(ctx context.Context, queuejob *QueueJob, path string) error {
	return nil
}

func (r *ResumeHandler) ParseResume(ctx context.Context, fileBytes []byte) (*llm.ResumeElements, error) {

	resumeText, err := ExtractPDFText(fileBytes)
	if err != nil {
		return nil, err
	}
	elements, err := r.resumeLLM.GetResumeElements(ctx, resumeText, llm.RESUME_EXTRACTION_PROMPT)
	if err != nil {
		return nil, err
	}

	return elements, nil

}


