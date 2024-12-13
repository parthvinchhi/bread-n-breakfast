package repository

import "github.com/parthvinchhi/bread-n-breakfast/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) error
}
