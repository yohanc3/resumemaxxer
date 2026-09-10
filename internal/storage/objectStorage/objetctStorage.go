package objectstorage

import "errors"

type ObjectStorage interface {
	Put(filename string) (string, error)
	Get(filename string) (string, error)
	Delete(filename string) (string, error)
}

type R2 struct {
	bucketName      string
	accessKeyID     string
	secretAccessKey string
	accountID       string
}

func (o R2) Put(filename string) (string, error) {
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
