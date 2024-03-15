package http

import (
	"github.com/DmitriyKomarovCoder/Film_library/internal/actor"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"net/http"
)

type Handler struct {
	usecase actor.Usecase
	log     logger.Logger
}

func NewHandler(usecase actor.Usecase, log logger.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		log:     log,
	}
}

// добавление информации об актёре (имя, пол, дата рождения),
func (h *Handler) AddActor(w http.ResponseWriter, r *http.Request) {

}

// изменение информации об актёре.
func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
}

// удаление информации об актёре,
func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
}

// получение списка актёров, для каждого актёра выдаётся также список фильмов с его участием,
func (h *Handler) GetActor(w http.ResponseWriter, r *http.Request) {
}
