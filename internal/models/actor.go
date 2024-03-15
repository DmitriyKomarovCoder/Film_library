package models

import "time"

type RequestActor struct {
	ActorID   uint      `json:"actor_id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type ResponseActor struct {
	ActorID   uint      `json:"actor_id"`
	Name      string    `json:"name,,omitempty"`
	Gender    string    `json:"gender,omitempty"`
	BirthDate time.Time `json:"birth_date,omitempty"`
	Movie     []uint    `json:"movie,omitempty"`
}
