package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/internal/movie"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type Handler struct {
	usecase movie.Usecase
	log     logger.Logger
}

func NewHandler(usecase movie.Usecase, log logger.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		log:     log,
	}
}

// добавление информации о фильме.
/*
При добавлении фильма указываются его название (не менее 1 и не более 150 символов),
описание (не более 1000 символов), дата выпуска, рейтинг (от 0 до 10) и список актёров:
*/

// @Summary Add a new movie
// @Tags movies
// @Description Add a new movie to the database
// @Accept json
// @Produce json
// @Param movie body models.CreateMovie true "Movie details"
// @Success 200 {object} uint "Movie ID"
// @Failure 400 {object} ResponseError "Bad Request"
// @Failure 500 {object} ResponseError "Internal Server Error"
// @Router /movie/add [post]
func (h *Handler) AddMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.CreateMovie
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.usecase.CreateMovie(&movie)
	if err != nil {
		if errors.Is(err, &models.ErrNilIDActor{}) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.log.Errorf("user bad request: %s", err)
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

// @Summary Update movie
// @Tags movies
// @Description Update movie to the database
// @Accept json
// @Produce json
// @Param movie body models.UpdateMovie true "Movie details"
// @Success 200 {object} models.UpdateMovie "Update movie response"
// @Failure 400 {object} ResponseError "Bad Request"
// @Failure 500 {object} ResponseError "Internal Server Error"
// @Router /movie/update [put]
func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	var updMovie models.UpdateMovie
	updMovie.Rating = -1

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&updMovie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(updMovie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resMovie, err := h.usecase.UpdateMovie(&updMovie)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "update object does not exist", http.StatusBadRequest)
			h.log.Errorf("user bad request: %s", err)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	responseJSON, err := json.Marshal(resMovie)
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

/*
получение списка фильмов с возможностью сортировки по названию, по рейтингу, по дате выпуска.
По умолчанию используется сортировка по рейтингу (по убыванию),
*/

// @Summary		Get sort movies
// @Tags		movies
// @Description	Get sort movies use query
// @Produce		json
// @Param   enumstring  query     string     false  "string enums"       Enums(A, B, C)
// @Param   enumint     query     int        false  "int enums"          Enums(1, 2, 3)
// @Param   enumnumber  query     number     false  "int enums"          Enums(1.1, 1.2, 1.3)
// @Param   string      query     string     false  "string valid"       minlength(5)  maxlength(10)
// @Param   int         query     int        false  "int valid"          minimum(1)    maximum(10)
// @Param   default     query     string     false  "string default"     default(A)
// @Param   example     query     string     false  "string example"     example(string)
// @Param   collection  query     []string   false  "string collection"  collectionFormat(multi)
// @Param   extensions  query     []string   false  "string collection"  extensions(x-example=test,x-nullable)
// @Success		200		{object}	[]models.ResponseMovie  "Response array"
// @Failure		400		{object}	ResponseError			 "Client error"
// @Failure		500		{object}	ResponseError			 "Internal Server Error"
// @Router		/film [get]
func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {
	var query, value string
	queryParams := r.URL.Query()

	ratingValue := queryParams.Get("rating")
	titleValue := queryParams.Get("title")
	releaseDateValue := queryParams.Get("release_date")

	if ratingValue == "desc" || ratingValue == "asc" {
		query = "rating"
		value = ratingValue
	} else if query == "" && titleValue == "desc" || titleValue == "asc" {
		query = "title"
		value = titleValue
	} else if query == "" && releaseDateValue == "desc" || releaseDateValue == "acs" {
		query = "realease_date"
		value = releaseDateValue
	} else if query == "" {
		query = "rating"
		value = "desc"
	}

	masMovie, err := h.usecase.GetMovies(query, value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "movies does not exist", http.StatusBadRequest)
			h.log.Errorf("user bad request: %s", err)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	responseJSON, err := json.Marshal(masMovie)
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

// поиск фильма по фрагменту названия, по фрагменту имени актёра
func (h *Handler) SearchMovie(w http.ResponseWriter, r *http.Request) {
}

// @Summary		Get sort movies
// @Tags		movies
// @Description	Get sort movies use query
// @Produce		json
// @Param       uint         query      uint        false  "uint valid"
// @Success		200		{object}	uint  					 "Delete id"
// @Failure		400		{object}	ResponseError			 "Client error"
// @Failure		500		{object}	ResponseError			 "Internal Server Error"
// @Router		/film/delete [delete]
func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	idMovie, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id not int", http.StatusBadRequest)
		h.log.Errorf("user bad request: %s", err)
		return
	}

	id, err := h.usecase.DeleteMovie(uint(idMovie))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "movies does not exist", http.StatusBadRequest)
			h.log.Errorf("user bad request: %s", err)
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
