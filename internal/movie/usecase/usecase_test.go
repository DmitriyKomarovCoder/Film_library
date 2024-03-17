package usecase

import (
	"errors"
	"io"
	"reflect"
	"testing"
	"time"

	au "github.com/DmitriyKomarovCoder/Film_library/internal/actor/mocks"
	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	mr "github.com/DmitriyKomarovCoder/Film_library/internal/movie/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestUsecase_CreateMovie(t *testing.T) {
	mockResponse := models.CreateMovie{Title: "Forest Gamp", Description: "...", Rating: 10, ReleaseDate: time.Now()}

	tests := []struct {
		name             string
		expectedErr      error
		expectedID       uint
		mockUsecaseFn    func(*au.MockUsecase)
		mockRepositoryFn func(*mr.MockRepository)
	}{
		{
			name:        "Successful Add movie",
			expectedErr: nil,
			expectedID:  uint(1),
			mockUsecaseFn: func(mockUsecase *au.MockUsecase) {
				mockUsecase.EXPECT().CheckActors(gomock.Any()).Return(true, nil)
			},
			mockRepositoryFn: func(MockRepository *mr.MockRepository) {
				MockRepository.EXPECT().CreateMovie(gomock.Any()).Return(uint(1), nil)
			},
		},
		{
			name:        "Nil actor error",
			expectedErr: &models.ErrNilIDActor{},
			expectedID:  uint(0),
			mockUsecaseFn: func(mockUsecase *au.MockUsecase) {
				mockUsecase.EXPECT().CheckActors(gomock.Any()).Return(false, nil)
			},
			mockRepositoryFn: func(MockRepository *mr.MockRepository) {
			},
		},
		{
			name:        "Nil actor error",
			expectedErr: errors.New("some err"),
			expectedID:  uint(0),
			mockUsecaseFn: func(mockUsecase *au.MockUsecase) {
				mockUsecase.EXPECT().CheckActors(gomock.Any()).Return(true, errors.New("some err"))
			},
			mockRepositoryFn: func(MockRepository *mr.MockRepository) {
			},
		},
		{
			name:        "return error in usecase function",
			expectedErr: errors.New("some err"),
			expectedID:  uint(0),
			mockUsecaseFn: func(mockUsecase *au.MockUsecase) {
				mockUsecase.EXPECT().CheckActors(gomock.Any()).Return(true, nil)
			},
			mockRepositoryFn: func(MockRepository *mr.MockRepository) {
				MockRepository.EXPECT().CreateMovie(gomock.Any()).Return(uint(0), errors.New("some err"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecaseA := au.NewMockUsecase(ctrl)
			mockRepositoryM := mr.NewMockRepository(ctrl)

			tt.mockUsecaseFn(mockUsecaseA)
			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM, mockUsecaseA)

			id, err := mockUsecase.CreateMovie(&mockResponse)

			assert.Equal(t, tt.expectedID, id)

			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) || (tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUsecase_UpdateMovie(t *testing.T) {
	mockResponse := models.UpdateMovie{Title: "Forest Gamp", Description: "...", Rating: 10, ReleaseDate: time.Now()}

	tests := []struct {
		name             string
		expectedErr      error
		expectedStruct   *models.UpdateMovie
		mockRepositoryFn func(*mr.MockRepository)
		mockRequest      io.Reader
	}{
		{
			name:           "Successful update movie",
			expectedErr:    nil,
			expectedStruct: &mockResponse,
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().GetMovie(gomock.Any()).Return(models.UpdateMovie{}, nil)
				mockRepository.EXPECT().UpdateMovie(gomock.Any()).Return(&mockResponse, nil)
			},
		},
		{
			name:           "Erorr get movie",
			expectedErr:    errors.New("some error"),
			expectedStruct: &models.UpdateMovie{},
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().GetMovie(gomock.Any()).Return(models.UpdateMovie{}, errors.New("some error"))
			},
		},
		{
			name:           "return error in usecase function",
			expectedErr:    errors.New("some error"),
			expectedStruct: &models.UpdateMovie{},
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().GetMovie(gomock.Any()).Return(models.UpdateMovie{}, nil)
				mockRepository.EXPECT().UpdateMovie(gomock.Any()).Return(&models.UpdateMovie{}, errors.New("some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepositoryM := mr.NewMockRepository(ctrl)

			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM, nil)

			reqM := &models.UpdateMovie{}
			reqM.Rating = -1

			uMovie, err := mockUsecase.UpdateMovie(reqM)

			assert.Equal(t, tt.expectedStruct, uMovie)

			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) || (tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUsecase_DeleteMovie(t *testing.T) {
	tests := []struct {
		name             string
		expectedErr      error
		expectedID       uint
		mockRepositoryFn func(*mr.MockRepository)
	}{
		{
			name:        "Successful deletion",
			expectedErr: nil,
			expectedID:  uint(1),
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().GetMovie(gomock.Any()).Return(models.UpdateMovie{}, nil)
				mockRepository.EXPECT().DeleteMovie(gomock.Any()).Return(uint(1), nil)
			},
		},
		{
			name:        "Error getting movie",
			expectedErr: errors.New("some error"),
			expectedID:  uint(0),
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().GetMovie(gomock.Any()).Return(models.UpdateMovie{}, errors.New("some error"))
			},
		},
		{
			name:        "Error deleting movie",
			expectedErr: errors.New("some error"),
			expectedID:  0,
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().GetMovie(gomock.Any()).Return(models.UpdateMovie{}, nil)
				mockRepository.EXPECT().DeleteMovie(gomock.Any()).Return(uint(0), errors.New("some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepositoryM := mr.NewMockRepository(ctrl)

			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM, nil)

			deletedID, err := mockUsecase.DeleteMovie(1)

			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) || (tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}

			if deletedID != tt.expectedID {
				t.Errorf("Expected ID: %d, but got: %d", tt.expectedID, deletedID)
			}
		})
	}
}

func TestUsecase_GetMovies(t *testing.T) {
	mockResponse := []models.ResponseMovie{
		{MovieID: 1, Title: "Movie 1", Rating: 8},
		{MovieID: 2, Title: "Movie 2", Rating: 7},
	}

	tests := []struct {
		name             string
		expectedErr      error
		expectedMovies   []models.ResponseMovie
		mockRepositoryFn func(*mr.MockRepository)
		querySort        string
		direction        string
	}{
		{
			name:           "Successful get movies",
			expectedErr:    nil,
			expectedMovies: mockResponse,
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(mockResponse, nil)
			},
			querySort: "title",
			direction: "asc",
		},
		{
			name:           "Error getting movies",
			expectedErr:    errors.New("some error"),
			expectedMovies: nil,
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().GetMovies(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			querySort: "title",
			direction: "asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepositoryM := mr.NewMockRepository(ctrl)

			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM, nil)

			movies, err := mockUsecase.GetMovies(tt.querySort, tt.direction)

			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) || (tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(tt.expectedMovies, movies) {
				t.Errorf("Expected movies: %v, but got: %v", tt.expectedMovies, movies)
			}
		})
	}
}

func TestUsecase_SearchMovie(t *testing.T) {
	mockResponse := []models.ResponseMovie{
		{MovieID: 1, Title: "Movie 1", Rating: 8},
	}

	tests := []struct {
		name             string
		expectedErr      error
		expectedMovies   []models.ResponseMovie
		mockRepositoryFn func(*mr.MockRepository)
		actorName        string
		filmName         string
	}{
		{
			name:           "Successful search movie",
			expectedErr:    nil,
			expectedMovies: mockResponse,
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().SearchMovie(gomock.Any(), gomock.Any()).Return(mockResponse, nil)
			},
			actorName: "Tom Hanks",
			filmName:  "Forest Gump",
		},
		{
			name:           "Error searching movie",
			expectedErr:    errors.New("some error"),
			expectedMovies: nil,
			mockRepositoryFn: func(mockRepository *mr.MockRepository) {
				mockRepository.EXPECT().SearchMovie(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			actorName: "Tom Hanks",
			filmName:  "Forest Gump",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepositoryM := mr.NewMockRepository(ctrl)

			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM, nil)

			movies, err := mockUsecase.SearchMovie(tt.actorName, tt.filmName)

			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) || (tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(tt.expectedMovies, movies) {
				t.Errorf("Expected movies: %v, but got: %v", tt.expectedMovies, movies)
			}
		})
	}
}
