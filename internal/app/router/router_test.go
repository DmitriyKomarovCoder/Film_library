package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	delActor "github.com/DmitriyKomarovCoder/Film_library/internal/actor/delivery/http"
	mockActor "github.com/DmitriyKomarovCoder/Film_library/internal/actor/mocks"
	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestActorHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger, _ := test.NewNullLogger()
	fakeLogger := logger.Logger{mockLogger}

	mockUActor := mockActor.NewMockUsecase(ctrl)
	mockUActor.EXPECT().GetActors().Return([]models.ResponseActor{}, nil)

	actorService := delActor.NewHandler(mockUActor, fakeLogger)
	router := NewRouter(nil, actorService, &fakeLogger)

	req, err := http.NewRequest("GET", "/actors", nil)
	req.AddCookie(&http.Cookie{Name: "role", Value: "admin"})

	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	(*router).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
