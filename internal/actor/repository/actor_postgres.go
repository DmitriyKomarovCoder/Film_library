package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/postgres"
	"github.com/pkg/errors"
)

const (
	checkActor  = `SELECT COUNT(*) FROM actors WHERE actor_id = $1`
	createActor = `INSERT INTO actors (name, gender, birth_date)
				   VALUES ($1, $2, $3)
				   RETURNING actor_id;`

	getActor = `
		SELECT actor_id, name, gender, birth_date
		FROM actors
		WHERE actor_id = $1;
	`

	getActors = `
	SELECT actor_id, name, gender, birth_date
	FROM actors;`

	updateActor = `
		UPDATE actors SET 
			name = $1, 
			gender = $2, 
			birth_date = $3 
		WHERE actor_id = $4;
	`

	joinActor = `
	SELECT ma.movie_id, m.title
	FROM movie_actors ma
	JOIN movies m ON m.movie_id = ma.movie_id
	WHERE ma.actor_id = $1;`
)

type repository struct {
	db postgres.DbConn
}

func NewRepository(db postgres.DbConn) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateActor(actor *models.CreateActor) (uint, error) {
	var actorID uint
	err := r.db.QueryRow(context.Background(), createActor, actor.Name, actor.Gender, actor.BirthDate).Scan(&actorID)
	if err != nil {
		return 0, err
	}
	return actorID, nil
}

func (r *repository) UpdateActor(actor *models.UpdateActor) (*models.UpdateActor, error) {
	_, err := r.db.Exec(context.Background(), updateActor,
		actor.Name,
		actor.Gender,
		actor.BirthDate,
		actor.ActorID)

	if err != nil {
		return &models.UpdateActor{}, err
	}

	return actor, nil
}

func (r *repository) DeleteActor(actorID uint) (uint, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `DELETE FROM movie_actors WHERE actor_id = $1`, actorID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM actors WHERE actor_id = $1`, actorID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	return actorID, nil
}

func (r *repository) GetActors() ([]models.ResponseActor, error) {
	var actorRes []models.ResponseActor
	rows, err := r.db.Query(context.Background(), getActors)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var actor models.ResponseActor
		if err := rows.Scan(
			&actor.ActorID,
			&actor.Name,
			&actor.Gender,
			&actor.BirthDate,
		); err != nil {
			return nil, err
		}

		rowsMovie, err := r.db.Query(context.Background(), joinActor, actor.ActorID)
		if err != nil {
			return nil, err
		}

		for rowsMovie.Next() {
			var movieInf models.MovieArr
			if err := rowsMovie.Scan(
				&movieInf.Id,
				&movieInf.Title,
			); err != nil {
				return nil, err
			}
			actor.Movie = append(actor.Movie, movieInf)
		}
		if err := rowsMovie.Err(); err != nil {
			return nil, err
		}
		actorRes = append(actorRes, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actorRes, nil
}

func (r *repository) GetActor(actorID uint) (models.UpdateActor, error) {
	oldActor := models.UpdateActor{}
	row := r.db.QueryRow(context.Background(), getActor, actorID)

	err := row.Scan(&oldActor.ActorID, &oldActor.Name, &oldActor.Gender, &oldActor.BirthDate)
	if err != nil {
		return models.UpdateActor{}, err
	}

	return oldActor, nil
}

func (r *repository) CheckActor(actorID uint) (bool, error) {
	var count int
	err := r.db.QueryRow(context.Background(), checkActor, actorID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error in function CheckActor() layer Repository: %v", err)
	}

	return count > 0, nil
}
