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

	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/internal/movie/mocks"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestHandler_GetMovie(t *testing.T) {
	MockResponse := []models.ResponseMovie{{MovieID: 1, Title: "Forest Gamp", Description: "...", Rating: 10}}
	tests := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
		queryHader    string
		queryValue    string
	}{
		{
			name:         "Successful call to Get Movie with null query",
			expectedCode: http.StatusOK,
			expectedBody: `[{"movie_id":1,"title":"Forest Gamp","description":"...","release_date":"0001-01-01T00:00:00Z","rating":10}]`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(MockResponse, nil)
			},
		},
		{
			name:         "Successful rating query",
			expectedCode: http.StatusOK,
			expectedBody: `[{"movie_id":1,"title":"Forest Gamp","description":"...","release_date":"0001-01-01T00:00:00Z","rating":10}]`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(MockResponse, nil)
			},
			queryHader: "rating",
			queryValue: "desc",
		},
		{
			name:         "Successful title query",
			expectedCode: http.StatusOK,
			expectedBody: `[{"movie_id":1,"title":"Forest Gamp","description":"...","release_date":"0001-01-01T00:00:00Z","rating":10}]`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(MockResponse, nil)
			},
			queryHader: "title",
			queryValue: "asc",
		},
		{
			name:         "Successful release_date query",
			expectedCode: http.StatusOK,
			expectedBody: `[{"movie_id":1,"title":"Forest Gamp","description":"...","release_date":"0001-01-01T00:00:00Z","rating":10}]`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(MockResponse, nil)
			},
			queryHader: "release_date",
			queryValue: "asc",
		},
		{
			name:         "Error response server",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `Internal server error`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return([]models.ResponseMovie{}, errors.New("some error"))
			},
			queryHader: "title",
			queryValue: "asc",
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

			q := req.URL.Query()
			q.Add(tt.queryHader, tt.queryValue)
			req.URL.RawQuery = q.Encode()

			recorder := httptest.NewRecorder()

			mockHandler.GetMovie(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_UpdateMovie(t *testing.T) {
	mockResponse := models.UpdateMovie{MovieID: 1, Title: "Forest Gamp", Description: "...", Rating: 10}
	MockData, _ := json.Marshal(mockResponse)
	failedResponse := models.UpdateMovie{Title: "Forest Gamp", Description: "...", Rating: 10, ReleaseDate: time.Time{}}
	FailedData, _ := json.Marshal(failedResponse)

	tests := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
		mockRequest   io.Reader
	}{
		{
			name:         "Successful Update movie",
			expectedCode: http.StatusOK,
			expectedBody: string(MockData),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().UpdateMovie(gomock.Any()).Return(&mockResponse, nil)
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
			expectedBody:  `Invalid input body`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {},
			mockRequest:   bytes.NewBuffer(FailedData),
		},
		{
			name:         "Invalid pgx no rows error",
			expectedCode: http.StatusBadRequest,
			expectedBody: "update object does not exist",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().UpdateMovie(gomock.Any()).Return(&models.UpdateMovie{}, pgx.ErrNoRows)
			},
			mockRequest: bytes.NewBuffer(MockData),
		},
		{
			name:         "Invalid internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().UpdateMovie(gomock.Any()).Return(&models.UpdateMovie{}, errors.New("some err"))
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

			mockHandler.UpdateMovie(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_CreateMovie(t *testing.T) {
	mockResponse := models.CreateMovie{Title: "Forest Gamp", Description: "...", Rating: 10, ReleaseDate: time.Now()}
	MockData, _ := json.Marshal(mockResponse)
	failedResponse := models.CreateMovie{Title: "", Description: "...", Rating: 10, ReleaseDate: time.Now()}
	FailedData, _ := json.Marshal(failedResponse)

	tests := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
		mockRequest   io.Reader
	}{
		{
			name:         "Successful Add movie",
			expectedCode: http.StatusOK,
			expectedBody: `1`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().CreateMovie(gomock.Any()).Return(uint(1), nil)
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
			name:         "Invalid actor nil error",
			expectedCode: http.StatusBadRequest,
			expectedBody: "id actor not found",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().CreateMovie(gomock.Any()).Return(uint(0), &models.ErrNilIDActor{})
			},
			mockRequest: bytes.NewBuffer(MockData),
		},
		{
			name:         "Invalid internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().CreateMovie(gomock.Any()).Return(uint(0), errors.New("some err"))
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

			mockHandler.AddMovie(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_SearchMovie(t *testing.T) {
	MockResponse := []models.ResponseMovie{{MovieID: 1, Title: "Forest Gamp", Description: "...", Rating: 10}}
	tests := []struct {
		name            string
		expectedCode    int
		expectedBody    string
		mockUsecaseFn   func(*mocks.MockUsecase)
		queryActor      string
		queryValueActor string
		queryTitle      string
		queryValueTitle string
	}{
		{
			name:         "Successful Search",
			expectedCode: http.StatusOK,
			expectedBody: `[{"movie_id":1,"title":"Forest Gamp","description":"...","release_date":"0001-01-01T00:00:00Z","rating":10}]`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().SearchMovie(gomock.Any(), gomock.Any()).Return(MockResponse, nil)
			},
			queryActor:      "actor_name",
			queryValueActor: "lol",
			queryTitle:      "movie_title",
			queryValueTitle: "kek",
		},
		{
			name:         "Error null query",
			expectedCode: http.StatusBadRequest,
			expectedBody: `quey doesn't exist`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
			},
		},
		{
			name:         "Error response server",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `Internal server error`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().SearchMovie(gomock.Any(), gomock.Any()).Return([]models.ResponseMovie{}, errors.New("some error"))
			},
			queryActor:      "actor_name",
			queryValueActor: "lol",
			queryTitle:      "movie_title",
			queryValueTitle: "kek",
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

			q := req.URL.Query()
			q.Add(tt.queryActor, tt.queryValueActor)
			q.Add(tt.queryTitle, tt.queryValueTitle)
			req.URL.RawQuery = q.Encode()

			recorder := httptest.NewRecorder()

			mockHandler.SearchMovie(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_DeleteMovie(t *testing.T) {
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
				mockUsecase.EXPECT().DeleteMovie(gomock.Any()).Return(uint(1), nil)
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
				mockUsecase.EXPECT().DeleteMovie(gomock.Any()).Return(uint(0), pgx.ErrNoRows)
			},
			queryID:    "id",
			queryValue: "1",
		},
		{
			name:         "Invalid internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error",
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().DeleteMovie(gomock.Any()).Return(uint(0), errors.New("some err"))
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

			req := httptest.NewRequest("DELETE", "/movie/delete", nil)

			q := req.URL.Query()
			q.Add(tt.queryID, tt.queryValue)
			req.URL.RawQuery = q.Encode()

			recorder := httptest.NewRecorder()

			mockHandler.DeleteMovie(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}
