package objectstorage

type ObjectStorage interface {
	Put(filename string) (string, error)
	Get(filename string) (string, error)
	Delete(filename string) (string, error)
}

type R2 struct {
	bucketName string
	apiKey     string
}

func (o R2) Put(filename string) (string, error) {
	return "", nil
}

func (o R2) Get(filename string) (string, error) {
	return "", nil
}

func (o R2) Delete(filename string) (string, error) {
	return "", nil
}

func NewR2(bucketName string, apiKey string) *R2 {
	return &R2{bucketName: bucketName, apiKey: apiKey}
}
