package resumebuilder

import (
	"context"
	"database/sql"

	"github.com/yohanc3/resumemaxxer/internal/storage/objectStorage"
)

type QueueJob struct {
	ID           string
	JobPostingID string
	UserID       string
	CompanyName  string
	Retries      int
	Status       string
	CreatedAt    int64
	FulfilledAt  int64
}

type ResumeBuilder struct {
	db      *sql.DB
	storage objectstorage.ObjectStorage 
}

func NewResumeBuilder(sem chan struct{}, db *sql.DB, storage objectstorage.ObjectStorage) ResumeBuilder {
	return ResumeBuilder{db: db, storage: storage}
}

func (r *ResumeBuilder) CreateResume(context context.Context, queuejob *QueueJob) error {
	return nil
}
