package repository

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateActor(t *testing.T) {
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
			errRows:    pgx.ErrNoRows,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			repo := NewRepository(mock)

			actor := &models.CreateActor{
				Name:      "John Doe",
				Gender:    "Male",
				BirthDate: time.Now(),
			}

			escapedQuery := regexp.QuoteMeta(createActor)
			mock.ExpectQuery(escapedQuery).
				WithArgs(actor.Name, actor.Gender, actor.BirthDate).
				WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(test.returnRows)).
				WillReturnError(test.errRows)

			actorID, err := repo.CreateActor(actor)

			if test.errRows != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.returnRows, actorID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateActor(t *testing.T) {
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

			actor := &models.UpdateActor{
				Name:      "Updated Name",
				Gender:    "Male",
				BirthDate: time.Now(),
				ActorID:   1,
			}

			escapedQuery := regexp.QuoteMeta(updateActor)
			mock.ExpectExec(escapedQuery).
				WithArgs(actor.Name, actor.Gender, actor.BirthDate, actor.ActorID).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1)).
				WillReturnError(test.errRows)

			updatedActor, err := repo.UpdateActor(actor)

			if test.errRows != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, actor, updatedActor)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteActor(t *testing.T) {
	tests := []struct {
		name               string
		actorID            uint
		errRowsMovieActors error
		flag               bool
		errRowsActors      error
		transaction        bool
	}{
		{
			name:               "Success",
			actorID:            uint(1),
			errRowsMovieActors: nil,
			errRowsActors:      nil,
			flag:               true,
			transaction:        true,
		},
		{
			name:               "Error Movie_Actor table",
			actorID:            uint(1),
			errRowsMovieActors: errors.New("mock error"),
			errRowsActors:      nil,
			flag:               false,
			transaction:        false,
		},
		{
			name:               "Error actor table",
			actorID:            uint(3),
			errRowsMovieActors: nil,
			errRowsActors:      errors.New("mock error"),
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

			mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM movie_actors WHERE actor_id = $1`)).
				WithArgs(test.actorID).
				WillReturnResult(pgxmock.NewResult("DELETE", 1)).
				WillReturnError(test.errRowsMovieActors)

			if test.flag {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM actors WHERE actor_id = $1`)).
					WithArgs(test.actorID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1)).
					WillReturnError(test.errRowsActors)
			}

			if test.transaction {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}
			actorID, err := repo.DeleteActor(test.actorID)

			if test.errRowsMovieActors != nil || test.errRowsActors != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.actorID, actorID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteActorBegin(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewRepository(mock)

	mock.ExpectBegin().WillReturnError(errors.New("err"))

	_, err := repo.DeleteActor(uint(1))

	assert.Equal(t, err, errors.New("err"))
}

func TestDeleteActorCommitTest(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewRepository(mock)

	mock.ExpectBegin()

	actorID := uint(1)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM movie_actors WHERE actor_id = $1`)).
		WithArgs(actorID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)).
		WillReturnError(nil)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM actors WHERE actor_id = $1`)).
		WithArgs(actorID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)).
		WillReturnError(nil)

	_, err := repo.DeleteActor(uint(1))

	mock.ExpectCommit().WillReturnError(errors.New("err"))

	assert.Equal(t, err, errors.New("all expectations were already fulfilled, call to Commit transaction was not expected"))
}

func TestGetActor(t *testing.T) {
	time := time.Now()
	movieArr := []models.MovieArr{{Id: uint(1), Title: "name"}}
	tests := []struct {
		name           string
		rows           *pgxmock.Rows
		rowsMovie      *pgxmock.Rows
		rowsMovieErr   error
		err            error
		rowsErr        error
		expected       []models.ResponseActor
		expectedLast   bool
		errTransaction bool
	}{
		{
			name: "ValidFeed",
			rows: pgxmock.NewRows([]string{
				"actor_id", "name", "gender", "birth_date",
			}).AddRow(
				uint(1), "Voronin", "M", time,
			),
			rowsMovie: pgxmock.NewRows([]string{
				"ma.movie_id",
				"m.title",
			}).AddRow(
				uint(1),
				"name",
			),

			err:            nil,
			rowsErr:        nil,
			rowsMovieErr:   nil,
			expected:       []models.ResponseActor{{ActorID: uint(1), Name: "Voronin", Gender: "M", BirthDate: time, Movie: movieArr}},
			expectedLast:   false,
			errTransaction: true,
		},
		{
			name: "Invalid Scan Movie",
			rows: pgxmock.NewRows([]string{
				"actor_id", "name", "gender", "birth_date",
			}).AddRow(
				uint(1), "Voronin", "M", time,
			),
			rowsMovie: pgxmock.NewRows([]string{
				"movie_id",
				"title",
			}).AddRow(
				"string",
				"name",
			),

			err:            fmt.Errorf(""),
			rowsErr:        nil,
			rowsMovieErr:   nil,
			expected:       nil,
			expectedLast:   false,
			errTransaction: true,
		},
		{
			name: "Invalid Scan actor",
			rows: pgxmock.NewRows([]string{
				"actor_id", "name", "gender", "birth_date",
			}).AddRow(
				"string", "Voronin", "M", time,
			),

			err:            fmt.Errorf("[repo] Scanning value error for column 'id': Scan: invalid UUID length: 3"),
			rowsErr:        nil,
			rowsMovieErr:   nil,
			expected:       nil,
			expectedLast:   false,
			errTransaction: false,
		},
		{
			name: "Invalid Scan Movie",
			rows: pgxmock.NewRows([]string{
				"actor_id", "name", "gender", "birth_date",
			}).AddRow(
				uint(1), "Voronin", "M", time,
			),
			rowsMovie: pgxmock.NewRows([]string{
				"movie_id",
				"title",
			}).AddRow(
				uint(1),
				"name",
			),

			err:            errors.New("err"),
			rowsErr:        nil,
			rowsMovieErr:   errors.New("err"),
			expected:       nil,
			expectedLast:   false,
			errTransaction: true,
		},
		{
			name: "Rows error actor",
			rows: pgxmock.NewRows([]string{
				"actor_id", "name", "gender", "birth_date",
			}).RowError(0, errors.New("err")),
			rowsMovie: pgxmock.NewRows([]string{
				"movie_id",
				"title",
			}),

			err:            fmt.Errorf("err"),
			rowsErr:        nil,
			rowsMovieErr:   nil,
			expected:       nil,
			expectedLast:   false,
			errTransaction: false,
		},
		{
			name: "Rows error movie",
			rows: pgxmock.NewRows([]string{
				"actor_id", "name", "gender", "birth_date",
			}).AddRow(
				uint(1), "Voronin", "M", time,
			),
			rowsMovie: pgxmock.NewRows([]string{
				"movie_id",
				"title",
			}).RowError(0, errors.New("err")),

			err:            fmt.Errorf("err"),
			rowsErr:        nil,
			rowsMovieErr:   nil,
			expected:       nil,
			expectedLast:   false,
			errTransaction: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := NewRepository(mock)

			escapedQuery := regexp.QuoteMeta(getActors)
			mock.ExpectQuery(escapedQuery).
				WillReturnRows(test.rows).
				WillReturnError(test.rowsErr)

			if test.errTransaction {
				escapedQueryCategory := regexp.QuoteMeta(joinActor)
				mock.ExpectQuery(escapedQueryCategory).
					WillReturnRows(test.rows).
					WillReturnError(test.rowsMovieErr)
			}
			transactions, err := repo.GetActors()

			if !reflect.DeepEqual(transactions, test.expected) {
				t.Errorf("Expected transactions: %v, but got: %v", test.expected, transactions)
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
