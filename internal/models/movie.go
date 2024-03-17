package models

import "time"

type UpdateMovie struct {
	MovieID     uint      `json:"movie_id" validate:"required"`
	Title       string    `json:"title" validate:"max=150"`
	Description string    `json:"description" validate:"max=1000"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      int       `json:"rating" validate:"min=-1,max=10"`
}

type CreateMovie struct {
	Title       string    `json:"title" validate:"required,min=1,max=150"`
	Description string    `json:"description" validate:"max=1000"`
	ReleaseDate time.Time `json:"release_date" validate:"required"`
	Rating      int       `json:"rating" validate:"min=0,max=10"`
	Actors      []uint    `json:"actors"`
}

type ResponseMovie struct {
	MovieID     uint      `json:"movie_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      int       `json:"rating"`
}

//type Movie struct {
//	MovieID     int       `json:"movie_id"`
//	Title       string    `json:"title"`
//	Description string    `json:"description"`
//	ReleaseDate time.Time `json:"release_date"`
//	Rating      float32   `json:"rating"`
//	Actor       []uint    `json:"actor"`
//}
