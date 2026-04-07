package minio_client

type MinIOOperations interface {
	UploadFile(objectName, localFilePath, contentType string) (string, error)
	ObjectExists(objectName string) (bool, error)
	GetObjectURL(objectName string) string
}

var _ MinIOOperations = (*MinIOClient)(nil)
