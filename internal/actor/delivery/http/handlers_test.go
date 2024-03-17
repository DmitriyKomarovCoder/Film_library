package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DmitriyKomarovCoder/Film_library/internal/actor/mocks"
	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestHandler_UpdateActor(t *testing.T) {
	mockResponse := models.UpdateActor{ActorID: 1, Name: "Forest Gamp", Gender: "W", BirthDate: time.Time{}}
	MockData, _ := json.Marshal(mockResponse)
	failedResponse := models.UpdateActor{Name: "Forest Gamp", Gender: "W", BirthDate: time.Time{}}
	FailedData, _ := json.Marshal(failedResponse)

	tests := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
		mockRequest   io.Reader
	}{
		{
			name:         "Successful Update actor",
			expectedCode: http.StatusOK,
			expectedBody: string(MockData),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().UpdateActor(gomock.Any()).Return(&mockResponse, nil)
			},
			mockRequest: bytes.NewBuffer(MockData),
		},
		{
			name:          "Invalid JSON body error",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `invalid character 'i' looking for beginning of value`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {},
			mockRequest:   strings.NewReader("invalid json data"),
		},
		{
			name:          "Invalid value error",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `Invaild input body`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {},
			mockRequest:   bytes.NewBuffer(FailedData),
		},
		{
			name:         "Invalid pgx no rows error",
			expectedCode: http.StatusBadRequest,
			expectedBody: "update object does not exist",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().UpdateActor(gomock.Any()).Return(&models.UpdateActor{}, pgx.ErrNoRows)
			},
			mockRequest: bytes.NewBuffer(MockData),
		},
		{
			name:         "Invalid internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().UpdateActor(gomock.Any()).Return(&models.UpdateActor{}, errors.New("some err"))
			},
			mockRequest: bytes.NewBuffer(MockData),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger, _ := test.NewNullLogger()
			fakeLogger := logger.Logger{mockLogger}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, fakeLogger)

			req := httptest.NewRequest("PUT", "/actor/update", tt.mockRequest)

			recorder := httptest.NewRecorder()

			mockHandler.UpdateActor(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_AddActor(t *testing.T) {
	mockResponse := models.CreateActor{Name: "Forest Gamp", Gender: "W", BirthDate: time.Time{}}
	MockData, _ := json.Marshal(mockResponse)
	failedResponse := models.CreateActor{Name: "Forest Gamp", Gender: "alien", BirthDate: time.Time{}}
	FailedData, _ := json.Marshal(failedResponse)

	tests := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
		mockRequest   io.Reader
	}{
		{
			name:         "Successful Add actor",
			expectedCode: http.StatusOK,
			expectedBody: `1`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().CreateActor(gomock.Any()).Return(uint(1), nil)
			},
			mockRequest: bytes.NewBuffer(MockData),
		},
		{
			name:          "Invalid JSON body error",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `invalid character 'i' looking for beginning of value`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {},
			mockRequest:   strings.NewReader("invalid json data"),
		},
		{
			name:          "Invalid value error",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `invalid input body`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {},
			mockRequest:   bytes.NewBuffer(FailedData),
		},
		{
			name:         "Invalid internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().CreateActor(gomock.Any()).Return(uint(0), errors.New("some err"))
			},
			mockRequest: bytes.NewBuffer(MockData),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger, _ := test.NewNullLogger()
			fakeLogger := logger.Logger{mockLogger}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, fakeLogger)

			req := httptest.NewRequest("POST", "/movies/add", tt.mockRequest)

			recorder := httptest.NewRecorder()

			mockHandler.AddActor(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_DeleteActor(t *testing.T) {
	tests := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
		queryID       string
		queryValue    string
	}{
		{
			name:         "Successful Delete",
			expectedCode: http.StatusOK,
			expectedBody: "1",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().DeleteActor(gomock.Any()).Return(uint(1), nil)
			},
			queryID:    "id",
			queryValue: "1",
		},
		{
			name:          "Invalid id",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `id not uint`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {},
			queryID:       "id",
			queryValue:    "string",
		},
		{
			name:          "Invalid id",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `id not negative`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {},
			queryID:       "id",
			queryValue:    "-1",
		},
		{
			name:         "Invalid pgx no rows error",
			expectedCode: http.StatusBadRequest,
			expectedBody: "movies does not exist",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().DeleteActor(gomock.Any()).Return(uint(0), pgx.ErrNoRows)
			},
			queryID:    "id",
			queryValue: "1",
		},
		{
			name:         "Invalid internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().DeleteActor(gomock.Any()).Return(uint(0), errors.New("some err"))
			},
			queryID:    "id",
			queryValue: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger, _ := test.NewNullLogger()
			fakeLogger := logger.Logger{mockLogger}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, fakeLogger)

			req := httptest.NewRequest("DELETE", "/actor/delete", nil)

			q := req.URL.Query()
			q.Add(tt.queryID, tt.queryValue)
			req.URL.RawQuery = q.Encode()

			recorder := httptest.NewRecorder()

			mockHandler.DeleteActor(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_GetMovie(t *testing.T) {
	tests := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Error response server",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `Internal server error`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GetActors().Return([]models.ResponseActor{}, errors.New("some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger, _ := test.NewNullLogger()
			fakeLogger := logger.Logger{mockLogger}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, fakeLogger)

			req := httptest.NewRequest("GET", "/movie", nil)

			recorder := httptest.NewRecorder()

			mockHandler.GetActors(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}
