package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func TestHandlerS3Failure(t *testing.T) {
	// Create a new stubber
	stubber := testtools.NewStubber()

	// Create the handler with stubbed S3 client
	handler := NewHandler(*stubber.SdkConfig)

	// Create a mock S3 event
	s3Event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "test-bucket",
					},
					Object: events.S3Object{
						Key: "test-key.txt",
					},
				},
			},
		},
	}

	// Create the JSON body from the S3 event
	body, err := json.Marshal(s3Event)
	if err != nil {
		t.Fatalf("failed to marshal S3 event: %v", err)
	}

	// Create a new request with the S3 event body
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	// Add the S3 GetObject stub that will return an error
	stubber.Add(testtools.Stub{
		OperationName: "GetObject",
		Input: &s3.GetObjectInput{
			Bucket: aws.String("test-bucket"),
			Key:    aws.String("test-key.txt"),
		},
		Error: &testtools.StubError{
			Err:           fmt.Errorf("simulated S3 error"),
			ContinueAfter: true,
		},
	})

	// Call the handler
	handler.HandleRequest(rr, req)

	// Verify all stubs were called
	if err := stubber.VerifyAllStubsCalled(); err != nil {
		t.Errorf("not all stubs were called: %v", err)
	}

	// Check the response
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expectedBody := "Failed to get S3 data\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestS3Success(t *testing.T) {
	// Create a new stubber
	stubber := testtools.NewStubber()

	// Create the handler with stubbed S3 client
	handler := NewHandler(*stubber.SdkConfig)

	// Create a mock S3 event
	s3Event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "test-bucket",
					},
					Object: events.S3Object{
						Key: "test-key.txt",
					},
				},
			},
		},
	}

	// Create the JSON body from the S3 event
	body, err := json.Marshal(s3Event)
	if err != nil {
		t.Fatalf("failed to marshal S3 event: %v", err)
	}

	// Create a new request with the S3 event body
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	// Add the S3 GetObject stub that will return an error
	stubber.Add(testtools.Stub{
		OperationName: "GetObject",
		Input: &s3.GetObjectInput{
			Bucket: aws.String("test-bucket"),
			Key:    aws.String("test-key.txt"),
		},
		Output: &s3.GetObjectOutput{},
	})

	// Call the handler
	handler.HandleRequest(rr, req)

	// Verify all stubs were called
	if err := stubber.VerifyAllStubsCalled(); err != nil {
		t.Errorf("not all stubs were called: %v", err)
	}

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
