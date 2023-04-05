package Server

import (
	"github.com/Wuchieh/IAQMNotificationSystem/Database"
	"github.com/Wuchieh/IAQMNotificationSystem/Line"
	"github.com/Wuchieh/IAQMNotificationSystem/math"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"sync"
)

type requestApi struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Aqi float64 `json:"aqi"`
}

func dangerAlerts(c *gin.Context) {
	var ra requestApi
	if err := c.Bind(&ra); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"status": false})
		return
	}

	if ra.Aqi < 150.4 {
		c.JSON(400, gin.H{"status": false, "data": ra})
		return
	}
	dangerAlertLogic(c, ra)
	c.JSON(200, ra)
}

func dangerAlertLogic(c *gin.Context, ra requestApi) error {
	locations, err := Database.GetAllLocations()
	if err != nil {
		return err
	}

	var locationMap = make(map[uuid.UUID][]Database.Location)
	for _, location := range locations {
		if _, ok := locationMap[location.UserID]; ok {
			locationMap[location.UserID] = append(locationMap[location.UserID], location)
		} else {
			locationMap[location.UserID] = []Database.Location{location}
		}
	}
	var wg sync.WaitGroup
	for userid, v := range locationMap {
		wg.Add(1)
		go func(userid uuid.UUID, v []Database.Location) {
			defer wg.Done()
			var dangerPoints []Database.Location
			for _, location := range v {
				lat1 := math.Radians(location.Location.P.X)
				lng1 := math.Radians(location.Location.P.Y)

				lat2 := math.Radians(ra.Lat)
				lng2 := math.Radians(ra.Lng)
				distance := math.CalculateDistance(lat1, lng1, lat2, lng2)
				if distance <= float64(location.Range)/1000 {
					dangerPoints = append(dangerPoints, location)
				}
			}

			if len(dangerPoints) >= 1 {
				msg := formatMsg(dangerPoints)
				id, err := Database.GetLineIdFromUserId(userid)
				if err != nil {
					log.Println(err)
					return
				}
				Line.SendMessage(id, msg)
			}
		}(userid, v)
	}
	wg.Wait()
	return nil
}

func formatMsg(points []Database.Location) string {
	message := "{{dangerAlertLabel}}\n"
	for i, point := range points {
		if i == 0 {
			message += point.NiceName
		} else {
			message += ", " + point.NiceName
		}
	}
	message += "\n{{dangerAlertEnd}}"
	return message
}
