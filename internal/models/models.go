package models

import "time"

// Reservation holds reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// Users is the user model
type Users struct {
	Id          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rooms is the room model
type Rooms struct {
	Id        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions is the restriction model
type Restrictions struct {
	Id              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservations is the reservation model
type Reservations struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomId    int
	Room      Rooms
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RoomRestrictions is the Room Restriction model
type RoomRestrictions struct {
	Id            int
	StartDate     time.Time
	EndDate       time.Time
	RoomId        int
	Room          Rooms
	ReservationId int
	Reservation   Reservations
	RestrictionId int
	Restriction   Restrictions
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
