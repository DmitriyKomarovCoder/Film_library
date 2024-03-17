package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/DmitriyKomarovCoder/Film_library/internal/actor"
	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4"
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

// @Summary Add a new actor
// @Tags actors
// @Description Add a new actor to the database
// @Accept json
// @Produce json
// @Param movie body models.CreateActor true "Movie details"
// @Success 200 {object} uint "Actor ID"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /actors/add [post]
func (h *Handler) AddActor(w http.ResponseWriter, r *http.Request) {
	cActor := models.CreateActor{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&cActor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(cActor); err != nil {
		http.Error(w, "invalid input body", http.StatusBadRequest)
		return
	}

	id, err := h.usecase.CreateActor(&cActor)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	responseJSON, err := json.Marshal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		h.log.Error(err)
	}
}

// @Summary Update actor
// @Tags actors
// @Description Update actor to the database
// @Accept json
// @Produce json
// @Param movie body models.UpdateActor true "Actor details"
// @Success 200 {object} models.UpdateActor "Update movie response"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /actors/update [put]
func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	var updActor models.UpdateActor
	updActor.Gender = "N"
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&updActor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(updActor); err != nil {
		http.Error(w, "Invaild input body", http.StatusBadRequest)
		return
	}

	resActor, err := h.usecase.UpdateActor(&updActor)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "update object does not exist", http.StatusBadRequest)
			h.log.Infof("user bad request: %s", err)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	responseJSON, err := json.Marshal(resActor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		h.log.Error(err)
	}
}

// @Summary		Delete actor
// @Tags		actors
// @Description	Delete actor to the database
// @Produce		json
// @Param       uint         query      uint        false  "uint valid"
// @Success		200		{object}	uint  					 "Delete id"
// @Failure		400		{object}	string			 "Client error"
// @Failure		500		{object}	string			 "Internal Server Error"
// @Router		/movies/delete [delete]
func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idActor, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id not uint", http.StatusBadRequest)
		h.log.Infof("user bad request: %s", err)
		return
	}

	if idActor < 0 {
		http.Error(w, "id not negative", http.StatusBadRequest)
		return
	}

	id, err := h.usecase.DeleteActor(uint(idActor))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "movies does not exist", http.StatusBadRequest)
			h.log.Infof("user bad request: %s", err)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	responseJSON, err := json.Marshal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		h.log.Error(err)
	}
}

// @Summary Get actors
// @Description Get by actors array with film information.
// @Tags movies
// @Produce json
// @Success 200 {object} []models.ResponseActor "Array actor with film information"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /actor [get]
func (h *Handler) GetActors(w http.ResponseWriter, r *http.Request) {
	actorRes, err := h.usecase.GetActors()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	responseJSON, err := json.Marshal(actorRes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		h.log.Error(err)
	}
}
