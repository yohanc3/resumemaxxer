package objectstorage

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
)

type ObjectStorage interface {
	Put(file []byte, filename string) (string, error)
	Get(filename string) (string, error)
	Delete(filename string) (string, error)
}

type R2 struct {
	bucketName      string
	accessKeyID     string
	secretAccessKey string
	accountID       string
}

// Uploads a given file, with the given filename (string),
func (o R2) Put(file *os.File, objectKey string) (string, error) {
	return "", errors.New("R2 Put not implemented")
}

func (o R2) Get(filename string) (string, error) {
	return "", errors.New("R2 Get not implemented")
}

func (o R2) Delete(filename string) (string, error) {
	return "", errors.New("R2 Delete not implemented")
}

func NewR2(bucketName, accessKeyID, secretAccessKey, accountID string) *R2 {
	return &R2{bucketName: bucketName, accessKeyID: accessKeyID, secretAccessKey: secretAccessKey, accountID: accountID}
}

// Given a user id, and user name, returns a filename and a unique objectkey
func CreateFileMetadata(userID string, userFirstName string, userLastName string) (string, string) {
	fileName := fmt.Sprintf("%s_%s_Resume.pdf", userFirstName, userLastName)
	objectKey := fmt.Sprintf("%s_%s_%s", "resumes", userID, rand.Text())
	return fileName, objectKey
}
