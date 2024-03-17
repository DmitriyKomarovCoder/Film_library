package models

import "time"

type CreateActor struct {
	Name      string    `json:"name" validate:"required"`
	Gender    string    `json:"gender" validate:"oneof=M W"`
	BirthDate time.Time `json:"birth_date"`
}

type UpdateActor struct {
	ActorID   uint      `json:"actor_id" validate:"required"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender" validate:"oneof=M W N"`
	BirthDate time.Time `json:"birth_date"`
}

type RequestActor struct {
	ActorID   uint      `json:"actor_id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type ResponseActor struct {
	ActorID   uint       `json:"actor_id"`
	Name      string     `json:"name"`
	Gender    string     `json:"gender"`
	BirthDate time.Time  `json:"birth_date"`
	Movie     []MovieArr `json:"movie"`
}

type MovieArr struct {
	Id    uint   `json:"movie_id"`
	Title string `json:"title"`
}
