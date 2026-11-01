package resumebuilderparser

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

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
	resumellm     llm.ResumeLLM
}

func NewResumeHandler(db *sql.DB, objectStorage objectstorage.ObjectStorage, llm llm.ResumeLLM) ResumeHandler {
	return ResumeHandler{db: db, objectStorage: objectStorage, resumellm: llm}
}

func (r *ResumeHandler) CreateResume(ctx context.Context, queuejob *QueueJob, path string) error {
	return nil
}

func (r *ResumeHandler) ParseAndSaveResume(ctx context.Context, userID string, file *os.File) error {

	resumeText, err := ExtractPDFText(file)
	if err != nil {
		return fmt.Errorf("error when extracting pdf text: %w", err)
	}

	resumeElements, err := r.resumellm.GetResumeElements(ctx, resumeText, llm.RESUME_EXTRACTION_PROMPT)
	if err != nil {
		return fmt.Errorf("error when getting resume pdf elements from llm: %w", err)
	}

	fileName, objectKey := objectstorage.CreateFileMetadata(userID, resumeElements.Header.FirstName, resumeElements.Header.LastName)
	_, err = r.objectStorage.Put(file, objectKey)
	if err != nil {
		return fmt.Errorf("error when uploading user's resume: %w", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error when starting transaction: %w", err)
	}
	defer tx.Rollback()

	resumeElementsJSON, err := json.Marshal(resumeElements)
	if err != nil {
		return fmt.Errorf("error when marshaling resume elements: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO parsed_resumes (user_id, parsed_resume) VALUES $1, $2",
		userID,
		resumeElementsJSON,
	)

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO resumes (user_id, filename, object_key) VALUES $1, $2, $3",
		userID,
		fileName,
		objectKey,
	)

	if err != nil {
		return fmt.Errorf("error whens storing parsed resume into the database: %w", err)
	}

	return nil
}
