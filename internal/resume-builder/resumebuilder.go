package resumebuilder

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/yohanc3/resumemaxxer/internal/storage/objectStorage"
)

type QueueJob struct {
	ID                    string
	JobPostingID          string
	JobPostingDescription string
	UserID                string
	Retries               int
	Status                string
}

type ResumeBuilder struct {
	db            *sql.DB
	storage       objectstorage.ObjectStorage
	fetchInterval time.Duration
}

func NewResumeBuilder(db *sql.DB, storage objectstorage.ObjectStorage) ResumeBuilder {
	return ResumeBuilder{db: db, storage: storage}
}

func (r *ResumeBuilder) CreateResume(ctx context.Context, queuejob *QueueJob) error {
	return errors.New("resume generation not implemented")
}
