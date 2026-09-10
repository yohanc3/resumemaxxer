package resumebuilder

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/yohanc3/resumemaxxer/internal/storage/objectStorage"
)

type QueueJob struct {
	ID           string
	JobPostingID string
	UserID       string
	CompanyName  string
	Retries      int
	Status       string
	CreatedAt    time.Time 
	FulfilledAt  sql.NullTime
}

type ResumeBuilder struct {
	db      *sql.DB
	storage objectstorage.ObjectStorage 
}

func NewResumeBuilder(db *sql.DB, storage objectstorage.ObjectStorage) ResumeBuilder {
	return ResumeBuilder{db: db, storage: storage}
}

func (r *ResumeBuilder) CreateResume(context context.Context, queuejob *QueueJob) error {
	return errors.New("resume generation not implemented") 
}
