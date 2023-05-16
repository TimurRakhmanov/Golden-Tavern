package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/RakhmanovTimur/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("invalid room id when trying to insert reservation")
	}
	return 1, nil
}

// InsertRoomRestrictions inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("invalid room id when trying to insert room restriction")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(roomID int, start, end time.Time) (bool, error) {
	layout := "2006-01-02"
	str := "2060-01-02"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	failDateTest := "2100-01-02"
	failT, err := time.Parse(layout, failDateTest)
	if err != nil {
		log.Println(err)
	}

	if failT == start {
		return false, errors.New("error querying database")
	}

	if start.After(t) {
		return false, nil
	}

	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	var room models.Room

	layout := "2006-01-02"
	failDate := "2099-01-01"
	t, err := time.Parse(layout, failDate)
	if err != nil {
		log.Println(err)
	}

	if start.After(t) {
		return rooms, nil
	}

	room.ID = 1
	rooms = append(rooms, room)

	return rooms, nil
}

// Gets a room by ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("can't find a room")
	}

	return room, nil

}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil

}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations returns a slice of new reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// GetReservationByID gets reservation by id
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var res models.Reservation
	return res, nil
}

// UpdateReservation updates a user in the database
func (m *testDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

// DeleteReservation deletes one reservation by id
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// Updates processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

// AllRooms returns all rooms in the database
func (m *testDBRepo) AllRooms() ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	return restrictions, nil
}
