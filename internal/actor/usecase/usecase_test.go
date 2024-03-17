package usecase

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/DmitriyKomarovCoder/Film_library/internal/actor/mocks"
	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestUsecase_UpdateActo(t *testing.T) {
	mockResponse := models.UpdateActor{Name: "Forest Gamp", Gender: "M", BirthDate: time.Now()}

	tests := []struct {
		name             string
		expectedErr      error
		expectedStruct   *models.UpdateActor
		mockRepositoryFn func(*mocks.MockRepository)
		mockRequest      io.Reader
	}{
		{
			name:           "Successful update actor",
			expectedErr:    nil,
			expectedStruct: &mockResponse,
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().GetActor(gomock.Any()).Return(models.UpdateActor{}, nil)
				mockRepository.EXPECT().UpdateActor(gomock.Any()).Return(&mockResponse, nil)
			},
		},
		{
			name:           "Erorr get actor",
			expectedErr:    errors.New("some error"),
			expectedStruct: &models.UpdateActor{},
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().GetActor(gomock.Any()).Return(models.UpdateActor{}, errors.New("some error"))
			},
		},
		{
			name:           "return error in get actor function",
			expectedErr:    errors.New("some error"),
			expectedStruct: &models.UpdateActor{},
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().GetActor(gomock.Any()).Return(models.UpdateActor{}, nil)
				mockRepository.EXPECT().UpdateActor(gomock.Any()).Return(&models.UpdateActor{}, errors.New("some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepositoryM := mocks.NewMockRepository(ctrl)

			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM)

			reqM := &models.UpdateActor{}

			reqM.Gender = "N"

			uMovie, err := mockUsecase.UpdateActor(reqM)

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
		mockRepositoryFn func(*mocks.MockRepository)
	}{
		{
			name:        "Successful deletion",
			expectedErr: nil,
			expectedID:  uint(1),
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().GetActor(gomock.Any()).Return(models.UpdateActor{}, nil)
				mockRepository.EXPECT().DeleteActor(gomock.Any()).Return(uint(1), nil)
			},
		},
		{
			name:        "Error getting actor",
			expectedErr: errors.New("some error"),
			expectedID:  uint(0),
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().GetActor(gomock.Any()).Return(models.UpdateActor{}, errors.New("some error"))
			},
		},
		{
			name:        "Error deleting actor",
			expectedErr: errors.New("some error"),
			expectedID:  0,
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().GetActor(gomock.Any()).Return(models.UpdateActor{}, nil)
				mockRepository.EXPECT().DeleteActor(gomock.Any()).Return(uint(0), errors.New("some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepositoryM := mocks.NewMockRepository(ctrl)

			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM)

			deletedID, err := mockUsecase.DeleteActor(1)

			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) || (tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}

			if deletedID != tt.expectedID {
				t.Errorf("Expected ID: %d, but got: %d", tt.expectedID, deletedID)
			}
		})
	}
}

func TestUsecase_CreateActor(t *testing.T) {
	mockResponse := uint(1)
	mockRequest := &models.CreateActor{
		Name:      "Liony Voronin",
		Gender:    "M",
		BirthDate: time.Now(),
	}

	tests := []struct {
		name             string
		expectedErr      error
		expectedID       uint
		mockRepositoryFn func(*mocks.MockRepository)
	}{
		{
			name:        "Successful actor creation",
			expectedErr: nil,
			expectedID:  mockResponse,
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().CreateActor(gomock.Any()).Return(mockResponse, nil)
			},
		},
		{
			name:        "Error creating actor",
			expectedErr: errors.New("some error"),
			expectedID:  0,
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().CreateActor(gomock.Any()).Return(uint(0), errors.New("some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepositoryM := mocks.NewMockRepository(ctrl)

			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM)

			createdID, err := mockUsecase.CreateActor(mockRequest)

			assert.Equal(t, tt.expectedID, createdID)

			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) || (tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetActors(t *testing.T) {
	mockResponse := []models.ResponseActor{
		{
			ActorID:   uint(1),
			Name:      "John Doe",
			Gender:    "M",
			BirthDate: time.Now(),
			Movie: []models.MovieArr{
				{Title: "Movie 1", Id: uint(1)},
				{Title: "Movie 2", Id: uint(2)},
			},
		},
	}

	tests := []struct {
		name             string
		expectedErr      error
		expectedActors   []models.ResponseActor
		mockRepositoryFn func(*mocks.MockRepository)
	}{
		{
			name:           "Successful actor retrieval",
			expectedErr:    nil,
			expectedActors: mockResponse,
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().GetActors().Return(mockResponse, nil)
			},
		},
		{
			name:           "Error retrieving actors",
			expectedErr:    errors.New("some error"),
			expectedActors: nil,
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().GetActors().Return(nil, errors.New("some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepositoryM := mocks.NewMockRepository(ctrl)

			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM)

			actors, err := mockUsecase.GetActors()

			assert.Equal(t, tt.expectedActors, actors)

			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) || (tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUsecase_CheckActors(t *testing.T) {
	actors := []uint{1, 2, 3}

	tests := []struct {
		name             string
		expectedErr      error
		expectedFlag     bool
		mockRepositoryFn func(*mocks.MockRepository)
	}{
		{
			name:         "All actors exist",
			expectedErr:  nil,
			expectedFlag: true,
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				for _, id := range actors {
					mockRepository.EXPECT().CheckActor(id).Return(true, nil)
				}
			},
		},
		{
			name:         "One actor does not exist",
			expectedErr:  nil,
			expectedFlag: false,
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().CheckActor(gomock.Any()).Return(true, nil).Times(2)
				mockRepository.EXPECT().CheckActor(gomock.Any()).Return(false, nil).Times(1)
			},
		},
		{
			name:         "Error checking actor existence",
			expectedErr:  errors.New("some error"),
			expectedFlag: false,
			mockRepositoryFn: func(mockRepository *mocks.MockRepository) {
				mockRepository.EXPECT().CheckActor(gomock.Any()).Return(false, errors.New("some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepositoryM := mocks.NewMockRepository(ctrl)

			tt.mockRepositoryFn(mockRepositoryM)

			mockUsecase := NewUsecase(mockRepositoryM)

			flag, err := mockUsecase.CheckActors(actors)

			assert.Equal(t, tt.expectedFlag, flag)

			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) || (tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}
		})
	}
}
