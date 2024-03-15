package usecase

import (
	"github.com/DmitriyKomarovCoder/Film_library/internal/actor"
	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
)

type Usecase struct {
	actorRepo actor.Repository
	log       logger.Logger
}

func NewUsecase(ar actor.Repository) *Usecase {
	return &Usecase{
		actorRepo: ar,
	}
}

func (r *Usecase) CreateActor(actor *models.RequestActor) (uint, error) {
	return 0, nil
}

func (r *Usecase) UpdateActor(actor *models.RequestActor) (*models.RequestActor, error) {
	return &models.RequestActor{}, nil
}

func (r *Usecase) DeleteActor(actorID uint) (uint, error) {
	return 0, nil
}

func (r *Usecase) GetActor(actorID uint) ([]models.ResponseActor, error) {
	return []models.ResponseActor{}, nil
}

func (r *Usecase) CheckActors(actors []uint) (bool, error) {
	flag := true
	var err error
	for _, id := range actors {
		if flag, err = r.actorRepo.CheckActor(id); err != nil || !flag {
			return flag, err
		}
	}
	return flag, nil
}
