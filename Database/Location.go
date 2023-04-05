package Database

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"time"
)

type Location struct {
	ID       uuid.UUID
	Location pgtype.Point // X lat Y lng
	NiceName string
	Time     time.Time
	UserID   uuid.UUID
	Range    int
	DeleteAt bool
}

func getLocationsFromUserID(userID uuid.UUID) ([]Location, error) {
	var Locations []Location
	//fmt.Println("SELECT id,nick_name,location,range FROM location WHERE user_id = $1 AND delete_at = false", userID)
	q, err := Db.Query("SELECT id,nick_name,location,range FROM location WHERE user_id = $1 AND delete_at = false", userID)
	if err != nil {
		return nil, err
	}
	for q.Next() {
		var l Location
		var locationBytes []byte
		err = q.Scan(&l.ID,
			&l.NiceName,
			&locationBytes,
			&l.Range)
		if err != nil {
			log.Println(err)
			continue
		}
		err = l.Location.UnmarshalJSON(locationBytes)
		if err != nil {
			log.Println(err)
			continue
		}
		Locations = append(Locations, l)
	}
	return Locations, nil
}

func GetAllLocations() ([]Location, error) {
	var Locations []Location
	//fmt.Println("SELECT id,nick_name,location,range FROM location WHERE user_id = $1 AND delete_at = false", userID)
	q, err := Db.Query("SELECT id,nick_name,location,range,user_id FROM location WHERE delete_at = false")
	if err != nil {
		return nil, err
	}
	for q.Next() {
		var l Location
		var locationBytes []byte
		err = q.Scan(&l.ID,
			&l.NiceName,
			&locationBytes,
			&l.Range,
			&l.UserID)
		if err != nil {
			log.Println(err)
			continue
		}
		err = l.Location.UnmarshalJSON(locationBytes)
		if err != nil {
			log.Println(err)
			continue
		}
		Locations = append(Locations, l)
	}
	return Locations, nil
}

func GetLineIdFromUserId(userid uuid.UUID) (string, error) {
	query, err := Db.Prepare(`SELECT "lineID" FROM users WHERE id = $1`)
	if err != nil {
		return "", err
	}
	defer query.Close()

	var id string

	err = query.QueryRow(userid).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
