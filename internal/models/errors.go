package models

type ErrNilIDActor struct{}

func (e *ErrNilIDActor) Error() string {
	return "id actor not found"
}
