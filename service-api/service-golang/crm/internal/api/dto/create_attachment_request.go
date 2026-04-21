package dto

type CreateAttachmentRequest struct {
	FileName      string `json:"fileName"`
	ContentType   string `json:"contentType"`
	StorageKey    string `json:"storageKey"`
	StorageDriver string `json:"storageDriver"`
	Source        string `json:"source"`
	UploadedBy    string `json:"uploadedBy"`
}
