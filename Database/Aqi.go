package Database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"time"
)

type Aqi struct {
	Id       int
	Location pgtype.Point
	Aqi      float64
	Time     time.Time
}

func GetAqis() (Aqis []Aqi) {
	q, err := Db.Query("SELECT * FROM aqi")
	if err != nil {
		log.Println(err)
		return nil
	}
	for q.Next() {
		var aqi Aqi
		var l []byte
		err = q.Scan(&aqi.Id,
			&l,
			&aqi.Aqi,
			&aqi.Time)
		if err != nil {
			log.Println(err)
		}
		if err = aqi.Location.UnmarshalJSON(l); err != nil {
			log.Println(err)
		}
		Aqis = append(Aqis, aqi)
	}
	return
}
