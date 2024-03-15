package repository

import (
	"context"

	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/postgres"
)

type Repository struct {
	*postgres.Postgres
}

const (
	createMovie = `
        INSERT INTO movies (title, description, release_date, rating)
        VALUES ($1, $2, $3, $4)
        RETURNING movie_id;
    `

	createMovieActor = `
		INSERT INTO movie_actors (movie_id, actor_id)
		VALUES ($1, $2)
    `

	updateMovie = `
		UPDATE movies SET title = $1, 
						  description = $2, 
						  release_date = $3, 
						  rating = $4 
		WHERE movie_id = $5;
	`

	getMovie = `
		SELECT movie_id, title, description, release_date, rating
		FROM movies
		ORDER BY
			CASE WHEN $1 = 'title'        AND $2 = 'asc'  THEN  title        END ASC,
			CASE WHEN $1 = 'title'        AND $2 = 'desc' THEN  title        END DESC,
			CASE WHEN $1 = 'rating'       AND $2 = 'asc'  THEN  title        END ASC,
			CASE WHEN $1 = 'rating'       AND $2 = 'desc' THEN  title        END DESC,
			CASE WHEN $1 = 'release_date' AND $2 = 'asc'  THEN  release_date END ASC,
			CASE WHEN $1 = 'release_date' AND $2 = 'desc' THEN  release_date END ASC,
	`
)

func New(pg *postgres.Postgres) *Repository {
	return &Repository{pg}
}

func (r *Repository) CreateMovie(movie *models.CreateMovie) (uint, error) {
	con, err := r.Pool.Acquire(context.Background())
	if err != nil {
		return 0, err
	}
	defer con.Release()

	tx, err := con.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	var movieID uint
	err = tx.QueryRow(context.Background(), createMovie, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).Scan(&movieID)
	if err != nil {
		return 0, err
	}

	for _, actorID := range movie.Actors {
		_, err = tx.Exec(context.Background(), createMovieActor, movieID, actorID)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	return movieID, nil
}

func (r *Repository) UpdateMovie(movie *models.UpdateMovie) (*models.UpdateMovie, error) {
	_, err := r.Pool.Exec(context.Background(), updateMovie,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.Rating,
		movie.MovieID)

	if err != nil {
		return &models.UpdateMovie{}, err
	}

	return movie, nil
}

func (r *Repository) DeleteMovie(movieID uint) (uint, error) {
	return 0, nil
}

func (r *Repository) GetMovies(querySort, direction string) ([]models.ResponseMovie, error) {
	rows, err := r.Pool.Query(context.Background(), getMovie, querySort, direction)
	if err != nil {
		return []models.ResponseMovie{}, err
	}
	defer rows.Close()

	var movies []models.ResponseMovie
	for rows.Next() {
		var movie models.ResponseMovie
		err := rows.Scan(&movie.MovieID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		if err != nil {
			return []models.ResponseMovie{}, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return []models.ResponseMovie{}, err
	}
	return movies, nil
}

func (r *Repository) GetMovie(movieID uint) (models.UpdateMovie, error) {
	return models.UpdateMovie{}, nil
}

func (r *Repository) SearchMovie(actorName, movieName string) (*models.ResponseMovie, error) {
	return &models.ResponseMovie{}, nil
}
