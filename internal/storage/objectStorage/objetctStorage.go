package objectstorage

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type ObjectStorage interface {
	Put(ctx context.Context, fileBytes []byte, objectKey string, mimeType string) (string, string, error)
	Get(filename string) (string, error)
	Delete(filename string) (string, error)
}

type R2 struct {
	r2Client        *s3.Client
	bucketName      string
	accessKeyID     string
	secretAccessKey string
	accountID       string
}

// Uploads a given file, with the given filename (string),
func (o R2) Put(ctx context.Context, fileBytes []byte, objectKey string, mimeType string) (string, string, error) {
	res, err := o.r2Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(o.bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(fileBytes),
		ContentType: aws.String(mimeType),
	})

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr != nil {
		return "", "", fmt.Errorf("error when uploading file")
	}

	return aws.ToString(res.ETag), aws.ToString(res.VersionId), nil
}

func (o R2) Get(filename string) (string, error) {
	return "", errors.New("R2 Get not implemented")
}

func (o R2) Delete(filename string) (string, error) {
	return "", errors.New("R2 Delete not implemented")
}

func NewR2(ctx context.Context, bucketName, accessKeyID, secretAccessKey, accountID string) (*R2, error) {
	r2Client, err := newR2Client(ctx, accountID, accessKeyID, secretAccessKey, bucketName)
	if err != nil {
		return nil, fmt.Errorf("error when creating r2 client: %w", err)
	}
	return &R2{r2Client: r2Client, bucketName: bucketName, accessKeyID: accessKeyID, secretAccessKey: secretAccessKey, accountID: accountID}, nil
}

func newR2Client(ctx context.Context, accountID, accessKey, secretKey, bucketName string) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		))
	if err != nil {
		return nil, fmt.Errorf("error when creating s3 default config: %w", err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(
			fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", accountID, bucketName),
		)
	}), nil
}

// Given a user id, and user name, returns a fileName and a unique objectKey
// The filename is formated as: "{Firstname}_{Lastname}_Resume"
func CreateFileMetadata(userID string, userFirstName string, userLastName string) (string, string) {
	fileName := fmt.Sprintf("%s_%s_Resume.pdf", userFirstName, userLastName)
	objectKey := fmt.Sprintf("%s_%s_%s", "resumes", userID, rand.Text())
	return fileName, objectKey
}
