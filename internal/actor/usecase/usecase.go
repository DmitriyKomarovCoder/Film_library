package usecase

import (
	"github.com/DmitriyKomarovCoder/Film_library/internal/actor"
	"github.com/DmitriyKomarovCoder/Film_library/internal/models"
)

type Usecase struct {
	actorRepo actor.Repository
}

func NewUsecase(ar actor.Repository) *Usecase {
	return &Usecase{
		actorRepo: ar,
	}
}

func (r *Usecase) CreateActor(actor *models.CreateActor) (uint, error) {
	return r.actorRepo.CreateActor(actor)
}

func (r *Usecase) UpdateActor(uActor *models.UpdateActor) (*models.UpdateActor, error) {
	oldActor, err := r.actorRepo.GetActor(uActor.ActorID)
	if err != nil {
		return &models.UpdateActor{}, err
	}

	if uActor.Name == "" {
		uActor.Name = oldActor.Name
	}

	if uActor.Gender == "N" {
		uActor.Gender = oldActor.Gender
	}

	if uActor.BirthDate.IsZero() {
		uActor.BirthDate = oldActor.BirthDate
	}

	return r.actorRepo.UpdateActor(uActor)
}

func (r *Usecase) DeleteActor(actorID uint) (uint, error) {
	if _, err := r.actorRepo.GetActor(actorID); err != nil {
		return 0, err
	}

	return r.actorRepo.DeleteActor(actorID)
}

func (r *Usecase) GetActors() ([]models.ResponseActor, error) {
	return r.actorRepo.GetActors()
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
