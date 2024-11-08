package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Handler struct {
	s3Client *s3.Client
}

func NewHandler(cfg aws.Config) *Handler {
	return &Handler{
		s3Client: s3.NewFromConfig(cfg),
	}
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {

	// decode the body into an s3 event
	var s3Event events.S3Event
	if err := json.NewDecoder(r.Body).Decode(&s3Event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get file content from S3
	_, err := h.s3Client.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: aws.String(s3Event.Records[0].S3.Bucket.Name),
		Key:    aws.String(s3Event.Records[0].S3.Object.Key),
	})

	if err != nil {
		http.Error(w, "Failed to get S3 data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
