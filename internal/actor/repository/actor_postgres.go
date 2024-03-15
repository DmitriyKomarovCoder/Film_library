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
	checkActor = `SELECT COUNT(*) FROM actors WHERE actor_id = $1`
)

type Repository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{pg}
}

func (r *Repository) CreateActor(actor *models.RequestActor) (uint, error) {
	return 0, nil
}

func (r *Repository) UpdateActor(actor *models.RequestActor) (*models.RequestActor, error) {
	return &models.RequestActor{}, nil
}

func (r *Repository) DeleteActor(actorID uint) (uint, error) {
	return 0, nil
}

func (r *Repository) GetActor(actorID uint) ([]models.ResponseActor, error) {
	return []models.ResponseActor{}, nil
}

func (r *Repository) CheckActor(actorID uint) (bool, error) {
	var count int
	err := r.Pool.QueryRow(context.Background(), checkActor, actorID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error in function CheckActor() layer Repository: %v", err)
	}

	return count > 0, nil
}
