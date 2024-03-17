package repository

import (
	"context"

	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

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

	getMovies = `
		SELECT movie_id, title, description, release_date, rating
		FROM movies
		ORDER BY
			CASE WHEN $1 = 'title' AND $2 = 'asc' THEN title END ASC,
			CASE WHEN $1 = 'title' AND $2 = 'desc' THEN title END DESC,
			CASE WHEN $1 = 'rating' AND $2 = 'asc' THEN rating END ASC,
			CASE WHEN $1 = 'rating' AND $2 = 'desc' THEN rating END DESC,
			CASE WHEN $1 = 'release_date' AND $2 = 'asc' THEN release_date END ASC,
			CASE WHEN $1 = 'release_date' AND $2 = 'desc' THEN release_date END DESC;
	`

	getMovie = `
		SELECT movie_id, title, description, release_date, rating
		FROM movies
		WHERE movie_id = $1;
	`
)

type repository struct {
	db postgres.DbConn
}

func NewRepository(db postgres.DbConn) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateMovie(movie *models.CreateMovie) (uint, error) {
	con, err := r.db.Acquire(context.Background())
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

func (r *repository) UpdateMovie(movie *models.UpdateMovie) (*models.UpdateMovie, error) {
	_, err := r.db.Exec(context.Background(), updateMovie,
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

func (r *repository) DeleteMovie(movieID uint) (uint, error) {
	con, err := r.db.Acquire(context.Background())
	if err != nil {
		return 0, err
	}
	defer con.Release()

	tx, err := con.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `DELETE FROM movie_actors WHERE movie_id = $1`, movieID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM movies WHERE movie_id = $1`, movieID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}
	return movieID, nil
}

func (r *repository) GetMovies(querySort, direction string) ([]models.ResponseMovie, error) {
	rows, err := r.db.Query(context.Background(), getMovies, querySort, direction)
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

func (r *repository) GetMovie(movieID uint) (models.UpdateMovie, error) {
	oldMovie := models.UpdateMovie{}
	row := r.db.QueryRow(context.Background(), getMovie, movieID)

	err := row.Scan(&oldMovie.MovieID, &oldMovie.Title, &oldMovie.Description, &oldMovie.ReleaseDate, &oldMovie.Rating)
	if err != nil {
		return models.UpdateMovie{}, err
	}

	return oldMovie, nil
}

func (r *repository) SearchMovie(actorName, movieName string) ([]models.ResponseMovie, error) {
	var movieArray []models.ResponseMovie
	var query string
	var row pgx.Rows
	var err error

	if actorName != "" && movieName != "" {
		query = `SELECT DISTINCT m.movie_id, m.title, m.description, m.release_date, m.rating
                FROM movies m
                JOIN movie_actors ma ON m.movie_id = ma.movie_id
                JOIN actors a ON ma.actor_id = a.actor_id
                WHERE a.name ILIKE '%' || $1 || '%' AND m.title ILIKE '%' || $2 || '%';`
		row, err = r.db.Query(context.Background(), query, actorName, movieName)
	} else if movieName != "" {
		query = `SELECT DISTINCT movie_id, title, description, release_date, rating
                FROM movies
                WHERE title ILIKE '%' || $1 || '%'`
		row, err = r.db.Query(context.Background(), query, movieName)
	} else if actorName != "" {
		query = `SELECT DISTINCT m.movie_id, m.title, m.description, m.release_date, m.rating
                FROM movies m
                JOIN movie_actors ma ON m.movie_id = ma.movie_id
                JOIN actors a ON ma.actor_id = a.actor_id
                WHERE a.name ILIKE '%' || $1 || '%';`
		row, err = r.db.Query(context.Background(), query, actorName)
	}

	if err != nil {
		return []models.ResponseMovie{}, err
	}
	defer row.Close()

	for row.Next() {
		var movie models.ResponseMovie
		err = row.Scan(&movie.MovieID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		if err != nil {
			return []models.ResponseMovie{}, err
		}
		movieArray = append(movieArray, movie)
	}

	if err = row.Err(); err != nil {
		return []models.ResponseMovie{}, err
	}

	return movieArray, nil
}
