package actor

import "github.com/DmitriyKomarovCoder/Film_library/internal/models"

type Usecase interface {
	CreateActor(actor *models.CreateActor) (uint, error)
	UpdateActor(actor *models.UpdateActor) (*models.UpdateActor, error)
	DeleteActor(actorID uint) (uint, error)
	GetActors() ([]models.ResponseActor, error)
	CheckActors(actors []uint) (bool, error)
}

type Repository interface {
	CreateActor(actor *models.CreateActor) (uint, error)
	UpdateActor(actor *models.UpdateActor) (*models.UpdateActor, error)
	DeleteActor(actorID uint) (uint, error)
	GetActor(actorID uint) (models.UpdateActor, error)
	GetActors() ([]models.ResponseActor, error)
	CheckActor(actor uint) (bool, error)
}
