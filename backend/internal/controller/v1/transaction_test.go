package v1

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andibalo/flip-test/internal/entity"
	"github.com/andibalo/flip-test/internal/model"
	"github.com/andibalo/flip-test/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionController_UploadCSV(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupRequest   func() (*http.Request, error)
		setupMock      func(mockService *mocks.MockTransactionService)
		expectedStatus int
	}{
		{
			name: "Success",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", "test.csv")
				if err != nil {
					return nil, err
				}
				part.Write([]byte("test content"))
				writer.Close()

				req, err := http.NewRequest(http.MethodPost, "/upload", body)
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			setupMock: func(mockService *mocks.MockTransactionService) {
				mockService.On("UploadCSVFile", mock.Anything, mock.Anything).Return(10, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing File",
			setupRequest: func() (*http.Request, error) {
				req, err := http.NewRequest(http.MethodPost, "/upload", nil)
				return req, err
			},
			setupMock:      func(mockService *mocks.MockTransactionService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid File Type",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", "test.txt")
				if err != nil {
					return nil, err
				}
				part.Write([]byte("test content"))
				writer.Close()

				req, err := http.NewRequest(http.MethodPost, "/upload", body)
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			setupMock:      func(mockService *mocks.MockTransactionService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Service Error",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", "test.csv")
				if err != nil {
					return nil, err
				}
				part.Write([]byte("test content"))
				writer.Close()

				req, err := http.NewRequest(http.MethodPost, "/upload", body)
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			setupMock: func(mockService *mocks.MockTransactionService) {
				mockService.On("UploadCSVFile", mock.Anything, mock.Anything).Return(0, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockTransactionService(t)
			tt.setupMock(mockService)

			controller := NewTransactionController(mockService)
			r := gin.Default()
			controller.AddRoutes(r)

			req, err := tt.setupRequest()
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestTransactionController_GetBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(mockService *mocks.MockTransactionService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			setupMock: func(mockService *mocks.MockTransactionService) {
				mockService.On("GetTotalBalance", mock.Anything).Return(int64(1000), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"total_balance":1000},"success":"success"}`,
		},
		{
			name: "Service Error",
			setupMock: func(mockService *mocks.MockTransactionService) {
				mockService.On("GetTotalBalance", mock.Anything).Return(int64(0), errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockTransactionService(t)
			tt.setupMock(mockService)

			controller := NewTransactionController(mockService)
			r := gin.Default()
			controller.AddRoutes(r)

			req, _ := http.NewRequest(http.MethodGet, "/balance", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestTransactionController_GetIssues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		setupMock      func(mockService *mocks.MockTransactionService)
		expectedStatus int
	}{
		{
			name:        "Success",
			queryParams: "?page=1&page_size=10",
			setupMock: func(mockService *mocks.MockTransactionService) {
				mockService.On("GetUnsuccessfulTransactions", mock.Anything, mock.Anything).Return([]*model.Transaction{}, int64(0), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Success with Sort",
			queryParams: "?page=1&page_size=10&sorts=+timestamp",
			setupMock: func(mockService *mocks.MockTransactionService) {
				mockService.On("GetUnsuccessfulTransactions", mock.Anything, mock.MatchedBy(func(filter entity.GetIssuesFilter) bool {
					return len(filter.Sorts.Data()) == 1 && filter.Sorts.Data()[0].Name == "timestamp"
				})).Return([]*model.Transaction{}, int64(0), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Sort Column",
			queryParams:    "?sorts=+invalid",
			setupMock:      func(mockService *mocks.MockTransactionService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Multiple Sort Columns",
			queryParams:    "?sorts=+timestamp,-status",
			setupMock:      func(mockService *mocks.MockTransactionService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Service Error",
			queryParams: "?page=1&page_size=10",
			setupMock: func(mockService *mocks.MockTransactionService) {
				mockService.On("GetUnsuccessfulTransactions", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockTransactionService(t)
			tt.setupMock(mockService)

			controller := NewTransactionController(mockService)
			r := gin.Default()
			controller.AddRoutes(r)

			req, _ := http.NewRequest(http.MethodGet, "/issues"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Logf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
