package usecase

import (
	"github.com/DmitriyKomarovCoder/Film_library/internal/actor"
	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/internal/movie"
)

type Usecase struct {
	mr movie.Repository
	au actor.Usecase
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

	if err != nil {
		return 0, err
	}
	return mu.mr.CreateMovie(film)
}

func (mu *Usecase) UpdateMovie(umovie *models.UpdateMovie) (*models.UpdateMovie, error) {
	oldMovie, err := mu.mr.GetMovie(umovie.MovieID)
	if err != nil {
		return &models.UpdateMovie{}, err
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
	if _, err := r.mr.GetMovie(filmID); err != nil {
		return 0, err
	}

	return r.mr.DeleteMovie(filmID)
}

func (r *Usecase) GetMovies(querySort, direction string) ([]models.ResponseMovie, error) {
	return r.mr.GetMovies(querySort, direction)
}

func (r *Usecase) SearchMovie(actorName, filmName string) ([]models.ResponseMovie, error) {
	return r.mr.SearchMovie(actorName, filmName)
}
