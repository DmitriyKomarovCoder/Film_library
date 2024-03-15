package usecase

import (
	"github.com/DmitriyKomarovCoder/Film_library/internal/actor"
	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/internal/movie"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
)

type Usecase struct {
	mr  movie.Repository
	au  actor.Usecase
	log logger.Logger
}

func NewUsecase(mr movie.Repository, au actor.Usecase) *Usecase {
	return &Usecase{
		mr: mr,
		au: au,
	}
}

func (mu *Usecase) CreateMovie(film *models.CreateMovie) (uint, error) {
	flag, err := mu.au.CheckActors(film.Actors)

	if !flag {
		return 0, &models.ErrNilIDActor{}
	}

	if err != nil || !flag {
		return 0, err
	}
	return mu.mr.CreateMovie(film)
}

func (mu *Usecase) UpdateMovie(umovie *models.UpdateMovie) (*models.UpdateMovie, error) {
	oldMovie, err := mu.mr.GetMovie(umovie.MovieID)
	if err != nil {
		return umovie, err
	}

	if umovie.Description == "" {
		umovie.Description = oldMovie.Description
	}

	if umovie.Rating == -1 {
		umovie.Rating = oldMovie.Rating
	}

	if umovie.Title == "" {
		umovie.Title = oldMovie.Title
	}

	if umovie.ReleaseDate.IsZero() {
		umovie.ReleaseDate = oldMovie.ReleaseDate
	}

	return mu.mr.UpdateMovie(umovie)
}

func (r *Usecase) DeleteMovie(filmID uint) (uint, error) {
	return 0, nil
}

func (r *Usecase) GetMovies(querySort, direction string) ([]models.ResponseMovie, error) {
	return r.mr.GetMovies(querySort, direction)
}

func (r *Usecase) SearchMovie(actorName, filmName string) (*models.ResponseMovie, error) {
	return &models.ResponseMovie{}, nil
}
