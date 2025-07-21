package aws

type S3Storage interface {
	CreatePresignedUploadURL(key string) (string, error)
}
