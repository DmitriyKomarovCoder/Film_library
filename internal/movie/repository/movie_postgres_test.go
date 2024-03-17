package repository

import (
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateMovie(t *testing.T) {
	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := NewRepository(mock)

	mockMovie := &models.CreateMovie{
		Title:       "Test Movie",
		Description: "Test Description",
		ReleaseDate: time.Now(),
		Rating:      8,
		Actors:      []uint{1},
	}

	expectedMovieID := uint(1)

	// 1 keys success

	mock.ExpectBegin()

	escapedCreateMovie := regexp.QuoteMeta(createMovie)

	mock.ExpectQuery(escapedCreateMovie).
		WithArgs(mockMovie.Title, mockMovie.Description, mockMovie.ReleaseDate, mockMovie.Rating).
		WillReturnRows(pgxmock.NewRows([]string{"movie_id"}).AddRow(expectedMovieID)).
		WillReturnError(nil)

	escapedCreateMovieActor := regexp.QuoteMeta(createMovieActor)
	mock.ExpectExec(escapedCreateMovieActor).
		WithArgs(expectedMovieID, uint(1)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1)).
		WillReturnError(nil)

	mock.ExpectCommit()

	movieID, err := repo.CreateMovie(mockMovie)

	assert.NoError(t, err)
	assert.Equal(t, expectedMovieID, movieID)

	// 1 keys success

	// 2 keys error in query createMovie

	mock.ExpectBegin()

	mock.ExpectQuery(escapedCreateMovie).
		WithArgs(mockMovie.Title, mockMovie.Description, mockMovie.ReleaseDate, mockMovie.Rating).
		WillReturnRows(pgxmock.NewRows([]string{"movie_id"})).
		WillReturnError(errors.New("err"))

	// escapedCreateMovieActor = regexp.QuoteMeta(createMovieActor)
	// // mock.ExpectExec(escapedCreateMovieActor).
	// 	WithArgs(expectedMovieID, uint(1)).
	// 	WillReturnResult(pgxmock.NewResult("INSERT", 1)).
	// 	WillReturnError(nil)

	mock.ExpectRollback()

	_, err = repo.CreateMovie(mockMovie)

	assert.Equal(t, err, errors.New("err"))
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// 2 keys error in query createMovie

	// 3 keys  error in exec movie actor

	mock.ExpectBegin()

	mock.ExpectQuery(escapedCreateMovie).
		WithArgs(mockMovie.Title, mockMovie.Description, mockMovie.ReleaseDate, mockMovie.Rating).
		WillReturnRows(pgxmock.NewRows([]string{"movie_id"}).AddRow(expectedMovieID)).
		WillReturnError(nil)

	mock.ExpectExec(escapedCreateMovieActor).
		WithArgs(expectedMovieID, uint(1)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1)).
		WillReturnError(errors.New("err"))

	mock.ExpectRollback()

	_, err = repo.CreateMovie(mockMovie)

	assert.Equal(t, err, errors.New("err"))

	// 3 keys  error in exec movie actor

	// 4 keys error in begin transaction

	mock.ExpectBegin().WillReturnError(errors.New("failed begin"))

	_, err = repo.CreateMovie(mockMovie)

	assert.Equal(t, err, errors.New("failed begin"))

	// 4 keys error in begin transaction

	// 5 keys error in done transaction

	mock.ExpectBegin()

	mock.ExpectQuery(escapedCreateMovie).
		WithArgs(mockMovie.Title, mockMovie.Description, mockMovie.ReleaseDate, mockMovie.Rating).
		WillReturnRows(pgxmock.NewRows([]string{"movie_id"}).AddRow(expectedMovieID)).
		WillReturnError(nil)

	mock.ExpectExec(escapedCreateMovieActor).
		WithArgs(expectedMovieID, uint(1)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1)).
		WillReturnError(nil)

	mock.ExpectCommit().WillReturnError(errors.New("err"))

	_, err = repo.CreateMovie(mockMovie)

	assert.Equal(t, err, errors.New("err"))

	// 5 keys error in done transaction

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateMovie(t *testing.T) {
	tests := []struct {
		name       string
		returnRows uint
		errRows    error
	}{
		{
			name:       "Success",
			returnRows: uint(1),
			errRows:    nil,
		},
		{
			name:       "Error",
			returnRows: uint(0),
			errRows:    errors.New("mock error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			repo := NewRepository(mock)

			movie := &models.UpdateMovie{
				Title:       "Updated Title",
				ReleaseDate: time.Now(),
				MovieID:     1,
				Description: "...",
				Rating:      8,
			}

			escapedQuery := regexp.QuoteMeta(updateMovie)
			mock.ExpectExec(escapedQuery).
				WithArgs(movie.Title, movie.Description, movie.ReleaseDate, movie.Rating, movie.MovieID).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1)).
				WillReturnError(test.errRows)

			updatedMovie, err := repo.UpdateMovie(movie)

			if test.errRows != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, movie, updatedMovie)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteMovie(t *testing.T) {
	tests := []struct {
		name               string
		movieID            uint
		errRowsMovieActors error
		flag               bool
		errRowsMovies      error
		transaction        bool
	}{
		{
			name:               "Success",
			movieID:            uint(1),
			errRowsMovieActors: nil,
			errRowsMovies:      nil,
			flag:               true,
			transaction:        true,
		},
		{
			name:               "Error Movie_Actor table",
			movieID:            uint(1),
			errRowsMovieActors: errors.New("mock error"),
			errRowsMovies:      nil,
			flag:               false,
			transaction:        false,
		},
		{
			name:               "Error movie table",
			movieID:            uint(3),
			errRowsMovieActors: nil,
			errRowsMovies:      errors.New("mock error"),
			flag:               true,
			transaction:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			repo := NewRepository(mock)

			mock.ExpectBegin()

			mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM movie_actors WHERE movie_id = $1`)).
				WithArgs(test.movieID).
				WillReturnResult(pgxmock.NewResult("DELETE", 1)).
				WillReturnError(test.errRowsMovieActors)

			if test.flag {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM movies WHERE movie_id = $1`)).
					WithArgs(test.movieID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1)).
					WillReturnError(test.errRowsMovies)
			}

			if test.transaction {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}
			actorID, err := repo.DeleteMovie(test.movieID)

			if test.errRowsMovieActors != nil || test.errRowsMovies != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.movieID, actorID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteMovieBegin(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewRepository(mock)

	mock.ExpectBegin().WillReturnError(errors.New("err"))

	_, err := repo.DeleteMovie(uint(1))

	assert.Equal(t, err, errors.New("err"))
}

func TestDeleteMovieCommitTest(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewRepository(mock)

	mock.ExpectBegin()

	actorID := uint(1)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM movie_actors WHERE movie_id = $1`)).
		WithArgs(actorID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)).
		WillReturnError(nil)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM movies WHERE movie_id = $1`)).
		WithArgs(actorID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)).
		WillReturnError(nil)

	_, err := repo.DeleteMovie(uint(1))

	mock.ExpectCommit().WillReturnError(errors.New("err"))

	assert.Equal(t, err, errors.New("all expectations were already fulfilled, call to Commit transaction was not expected"))
}

func TestRepository_GetMovieValid(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewRepository(mock)

	movieID := uint(1)
	expectedOldMovie := models.UpdateMovie{
		MovieID:     movieID,
		Title:       "Movie1",
		Description: "...",
		Rating:      8,
		ReleaseDate: time.Now(),
	}

	escapedQuery := regexp.QuoteMeta(getMovie)
	mock.ExpectQuery(escapedQuery).
		WithArgs(movieID).
		WillReturnRows(pgxmock.NewRows([]string{"movie_id", "title", "description", "release_date", "rating"}).
			AddRow(expectedOldMovie.MovieID, expectedOldMovie.Title, expectedOldMovie.Description, expectedOldMovie.ReleaseDate, expectedOldMovie.Rating))

	oldActor, err := repo.GetMovie(movieID)

	assert.NoError(t, err)
	assert.Equal(t, expectedOldMovie, oldActor)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_GetMovieNotValid(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewRepository(mock)

	actorID := uint(1)

	escapedQuery := regexp.QuoteMeta(getMovie)
	mock.ExpectQuery(escapedQuery).
		WithArgs(actorID).
		WillReturnRows(pgxmock.NewRows([]string{"movie_id", "title", "description", "release_date, rating"})).
		WillReturnError(errors.New("err"))
	_, err := repo.GetMovie(actorID)

	assert.Equal(t, err, errors.New("err"))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetMovies(t *testing.T) {
	successMovie := []models.ResponseMovie{{MovieID: uint(1), Title: "Voronin", Description: "...", ReleaseDate: time.Now(), Rating: 10}}
	tests := []struct {
		name     string
		rows     *pgxmock.Rows
		err      error
		rowsErr  error
		expected []models.ResponseMovie
	}{
		{
			name: "Success",
			rows: pgxmock.NewRows([]string{
				"movie_id", "title", "description", "release_date", "rating",
			}).AddRow(
				successMovie[0].MovieID, successMovie[0].Title, successMovie[0].Description, successMovie[0].ReleaseDate, successMovie[0].Rating,
			),
			err:      nil,
			rowsErr:  nil,
			expected: successMovie,
		},
		{
			name: "Failed query",
			rows: pgxmock.NewRows([]string{
				"movie_id", "title", "description", "release_date", "rating",
			}),
			err:      errors.New("err"),
			rowsErr:  errors.New("err"),
			expected: []models.ResponseMovie{},
		},
		{
			name: "Failed rows",
			rows: pgxmock.NewRows([]string{
				"movie_id", "title", "description", "release_date", "rating",
			}).RowError(0, errors.New("error rows")),
			err:      errors.New("error rows"),
			rowsErr:  nil,
			expected: []models.ResponseMovie{},
		},
		{
			name: "Failed in row 1",
			rows: pgxmock.NewRows([]string{
				"movie_id", "title", "description", "release_date", "rating",
			}).AddRow(successMovie[0].MovieID, successMovie[0].Title, successMovie[0].Description, successMovie[0].ReleaseDate, successMovie[0].Rating).
				RowError(1, errors.New("error rows")),
			err:      errors.New("error rows"),
			rowsErr:  nil,
			expected: []models.ResponseMovie{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := NewRepository(mock)

			escapedQuery := regexp.QuoteMeta(getMovies)
			mock.ExpectQuery(escapedQuery).
				WillReturnRows(test.rows).
				WillReturnError(test.rowsErr)

			actorsArray, err := repo.GetMovies("mock", "mock")

			if !reflect.DeepEqual(actorsArray, test.expected) {
				t.Errorf("Expected actors: %v, but got: %v", test.expected, actorsArray)
			}

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_SearchMovieSuccess(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewRepository(mock)

	successMovie := []models.ResponseMovie{{
		MovieID:     uint(1),
		Title:       "Voronin",
		Description: "...",
		ReleaseDate: time.Now(),
		Rating:      10}}

	// check with 2 paramerts

	escapedQuery := regexp.QuoteMeta(searchMovieActorMovieN)
	mock.ExpectQuery(escapedQuery).
		WithArgs("mock", "mock").
		WillReturnRows(pgxmock.NewRows([]string{"movie_id", "title", "description", "release_date", "rating"}).
			AddRow(successMovie[0].MovieID, successMovie[0].Title, successMovie[0].Description, successMovie[0].ReleaseDate, successMovie[0].Rating))

	oldActor, err := repo.SearchMovie("mock", "mock")

	assert.NoError(t, err)
	assert.Equal(t, successMovie, oldActor)

	// check with 1 parametrs movie

	escapedQuery = regexp.QuoteMeta(searchMovieName)
	mock.ExpectQuery(escapedQuery).
		WithArgs("mock").
		WillReturnRows(pgxmock.NewRows([]string{"movie_id", "title", "description", "release_date", "rating"}).
			AddRow(successMovie[0].MovieID, successMovie[0].Title, successMovie[0].Description, successMovie[0].ReleaseDate, successMovie[0].Rating))

	oldActor, err = repo.SearchMovie("", "mock")

	assert.NoError(t, err)
	assert.Equal(t, successMovie, oldActor)

	// check with 1 paramerts actor

	escapedQuery = regexp.QuoteMeta(searchMovieActorName)
	mock.ExpectQuery(escapedQuery).
		WithArgs("mock").
		WillReturnRows(pgxmock.NewRows([]string{"movie_id", "title", "description", "release_date", "rating"}).
			AddRow(successMovie[0].MovieID, successMovie[0].Title, successMovie[0].Description, successMovie[0].ReleaseDate, successMovie[0].Rating))

	oldActor, err = repo.SearchMovie("mock", "")

	assert.NoError(t, err)
	assert.Equal(t, successMovie, oldActor)

	// check with null parametrs

	oldActor, err = repo.SearchMovie("", "")

	assert.NoError(t, err)
	assert.Equal(t, []models.ResponseMovie{}, oldActor)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
