package Line

import (
	"fmt"
	"github.com/Wuchieh/IAQMNotificationSystem/Database"
	"github.com/Wuchieh/IAQMNotificationSystem/math"
	"log"
	"time"
)

var (
	aqis []Database.Aqi
)

func NotificationCronjob() {
	now := time.Now()

	for {
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		sleepDuration := nextHour.Sub(now)
		now = time.Now()
		fmt.Println("剩下", sleepDuration)
		time.Sleep(sleepDuration)

		notification()
	}
}

func notification() {
	nowHour := time.Now().Hour()
	//nowHour := 16
	var timeLimit []int
	for _, v := range notificationRange {
		if nowHour%v == 0 {
			timeLimit = append(timeLimit, v)
		}
	}
	sendNotification(timeLimit)
}

func sendNotification(limit []int) {
	if len(limit) < 1 {
		return
	}

	aqis = Database.GetAqis()

	if users, err := Database.GetUsersFromNoticeRange(limit); err != nil {
		log.Println(err)
		return
	} else {
		for _, u := range users {
			go func(u Database.User) {
				locations, err := u.SendNotification()
				if err != nil {
					log.Println(err)
				}
				msg := formatMsg(locations)
				if len(msg) >= 1 {
					SendMessage(u.LineId, msg)
				}
			}(u)
		}
	}
}

func formatMsg(locations []Database.Location) string {
	var msgs []string
	for _, location := range locations {
		var aqiAqi, count float64
		for _, aqi := range aqis {
			lat1 := math.Radians(location.Location.P.X)
			lng1 := math.Radians(location.Location.P.Y)

			lat2 := math.Radians(aqi.Location.P.X)
			lng2 := math.Radians(aqi.Location.P.Y)
			distance := math.CalculateDistance(lat1, lng1, lat2, lng2)
			if distance <= float64(location.Range)/1000 {
				aqiAqi += aqi.Aqi
				count++
			}
		}
		msgs = append(msgs, fmt.Sprintf("%s: %s", location.NiceName, sqiStatus(aqiAqi, count)))
	}
	if len(msgs) < 1 {
		return ""
	}
	message := "{{lineNotificationLabel}}\n"
	for _, msg := range msgs {
		message += fmt.Sprintf("\n%s", msg)
	}
	return message
}

func sqiStatus(aqi, count float64) string {
	if count == 0 {
		return "{{notFountData}}"
	}
	avg := aqi / count
	if avg >= 150.4 {
		return "{{danger}}"
	} else if 54.4 <= avg && avg < 150.4 {
		return "{{veryBad}}"
	} else if 35.4 <= avg && avg < 54.4 {
		return "{{bad}}"
	} else if 15.5 <= avg && avg < 35.4 {
		return "{{normal}}"
	} else if 0 <= avg && avg < 15.5 {
		return "{{good}}"
	}
	return "{{dataError}}"
}
