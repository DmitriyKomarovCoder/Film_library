package movie

import "github.com/DmitriyKomarovCoder/Film_library/internal/models"

type Usecase interface {
	CreateMovie(film *models.CreateMovie) (uint, error)
	UpdateMovie(film *models.UpdateMovie) (*models.UpdateMovie, error)
	DeleteMovie(filmID uint) (uint, error)
	GetMovies(querySort, direction string) ([]models.ResponseMovie, error)
	SearchMovie(actorName, filmName string) (*models.ResponseMovie, error)
}

type Repository interface {
	CreateMovie(movie *models.CreateMovie) (uint, error)
	UpdateMovie(movie *models.UpdateMovie) (*models.UpdateMovie, error)
	DeleteMovie(movieID uint) (uint, error)
	GetMovies(querySort, direction string) ([]models.ResponseMovie, error)
	SearchMovie(actorName, movieName string) (*models.ResponseMovie, error)
	GetMovie(movieID uint) (models.UpdateMovie, error)
}
