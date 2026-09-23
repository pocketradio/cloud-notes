package endpoints

import (
	"encoding/json"
	"io"
	"net/http"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileHandler struct {
	client *s3.Client
	bucket string
}

type File struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}

func NewFileHandler(client *s3.Client, bucket string) *FileHandler {
	return &FileHandler{
		client: client,
		bucket: bucket,
	}
}

func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	result, err := h.client.ListObjectsV2(r.Context(), &s3.ListObjectsV2Input{
		Bucket: aws.String(h.bucket),
	})
	if err != nil {
		http.Error(w, "failed to list files", http.StatusInternalServerError)
		return
	}

	files := make([]File, 0, len(result.Contents))
	for _, object := range result.Contents {
		files = append(files, File{
			Key:  aws.ToString(object.Key),
			Size: aws.ToInt64(object.Size),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(files)
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MiB max
		http.Error(w, "file is too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file") // returns the 1st file for the key
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// to get only the filename. like data.csv
	key := path.Base(header.Filename)
	if key == "." || key == "" {
		http.Error(w, "invalid file name", http.StatusBadRequest)
		return
	}

	_, err = h.client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(h.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})

	if err != nil {
		http.Error(w, "failed to upload file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"key": key,
	})
}

func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "file key is required", http.StatusBadRequest)
		return
	}

	result, err := h.client.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: aws.String(h.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		http.Error(w, "failed to get file", http.StatusNotFound)
		return
	}
	defer result.Body.Close()

	if result.ContentType != nil {
		w.Header().Set("Content-Type", aws.ToString(result.ContentType))
	}
	w.Header().Set("Content-Disposition", "inline")

	if _, err := io.Copy(w, result.Body); err != nil {
		return
	}
}
