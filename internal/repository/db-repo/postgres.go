package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/parthvinchhi/bread-n-breakfast/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int

	stmt := `insert into reservations (first_name,last_name,email,phone,start_date,
		end_date,room_id,created_at,updated_at) 
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date,end_date,room_id, reservation_id,
		restriction_id,created_at,updated_at) values($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomId,
		r.ReservationId,
		r.RestrictionId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomId return true if availability exists for roomId and false if no availability exists
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select 
			count(id) 
		from 
			room_restrictions 
		where
			room_id = $1 
			$2 < end_date and $3 > start_date`

	var numRow int

	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	if err := row.Scan(&numRow); err != nil {
		return false, err
	}
	if numRow == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		select
			r.id, r.room_name
		from
			rooms r
		where r.id not in 
		(select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);
		`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.Id,
			&room.RoomName,
		)

		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// GetRoomById gets a room by id
func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select * from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.Id, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}

// GetUserById returns a User by Id
func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
			from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.Id,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates the user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set 
		first_name = $1, 
		last_name = $2, 
		email = $3, 
		access_level = $4, 
		updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now())

	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates an user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", nil
	}

	return id, hashedPassword, nil
}

// AllReservations returns a slice of all the reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, 
			r.room_id, r.created_at, r.updated_at, r.processed,
			rm.id, rm.room_name
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			order by r.start_date asc
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.Id,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.Id,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// AllNewReservations returns a slice of all the reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, 
			r.room_id, r.created_at, r.updated_at,
			rm.id, rm.room_name
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			where processed = 0
			order by r.start_date asc
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.Id,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.Id,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// GetReservationById return one reservation by Id
func (m *postgresDBRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, 
			r.room_id, r.created_at, r.updated_at, r.processed,
			rm.id, rm.room_name
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			where r.id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomId,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.Id,
		&res.Room.RoomName,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

// UpdateReservations updates the user in the database
func (m *postgresDBRepo) UpdateReservations(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set 
		first_name = $1, 
		last_name = $2, 
		email = $3,
		phone = $4, 
		updated_at = $5
		where id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.Id)

	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation deletes one reservation by id
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateProcessedForReservation updates processed for reservation by id
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed = $1 where id = $2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select id,room_name,created_at,updated_at from rooms order by room_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}

	defer rows.Close()

	for rows.Next() {
		var rm models.Room
		err := rows.Scan(
			&rm.Id,
			&rm.RoomName,
			&rm.CreatedAt,
			&rm.UpdatedAt,
		)

		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, rm)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRestrictionsForRoomByDate returns retrictions for a room by date range
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restriction []models.RoomRestriction

	query := `
		select id, coalesce(reservation_id,0), restriction_id, room_id, start_date, end_date
		from room_restrictions where $1 < end_date and $2 >= start_date
		and room_id =$3
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.Id,
			&r.ReservationId,
			&r.RestrictionId,
			&r.RoomId,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return nil, err
		}
		restriction = append(restriction, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return restriction, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *postgresDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at)
		values ($1,$2,$3,$4,$5,$6)`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DeleteBlockById deletes a room restriction
func (m *postgresDBRepo) DeleteBlockById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from room_restrictions where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
