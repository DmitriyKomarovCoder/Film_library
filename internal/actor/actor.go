package actor

import "github.com/DmitriyKomarovCoder/Film_library/internal/models"

type Usecase interface {
	CreateActor(actor *models.RequestActor) (uint, error)
	UpdateActor(actor *models.RequestActor) (*models.RequestActor, error)
	DeleteActor(actorID uint) (uint, error)
	GetActor(actorID uint) ([]models.ResponseActor, error)
	CheckActors(actors []uint) (bool, error)
}

type Repository interface {
	CreateActor(actor *models.RequestActor) (uint, error)
	UpdateActor(actor *models.RequestActor) (*models.RequestActor, error)
	DeleteActor(actorID uint) (uint, error)
	GetActor(actorID uint) ([]models.ResponseActor, error)
	CheckActor(actor uint) (bool, error)
}
