package models

import "time"

// User is the user model
type User struct {
	Id          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the room model
type Room struct {
	Id        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the restriction model
type Restriction struct {
	Id              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the reservation model
type Reservation struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomId    int
	Room      Room
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RoomRestriction is the Room Restriction model
type RoomRestriction struct {
	Id            int
	StartDate     time.Time
	EndDate       time.Time
	RoomId        int
	Room          Room
	ReservationId int
	Reservation   Reservation
	RestrictionId int
	Restriction   Restriction
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
