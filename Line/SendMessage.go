package Line

import (
	"github.com/Wuchieh/IAQMNotificationSystem/i18n"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"strings"
)

func SendMessage(userId, message string) {
	var lang string
	p, err := bot.GetProfile(userId).Do()
	if err != nil {
		log.Println(err)
		lang = ""
	} else {
		lang = p.Language
	}
	message = strings.ReplaceAll(message, "{{lineNotificationLabel}}", i18n.Get(lang, "lineNotificationLabel"))
	message = strings.ReplaceAll(message, "{{notFountData}}", i18n.Get(lang, "notFountData"))
	message = strings.ReplaceAll(message, "{{danger}}", i18n.Get(lang, "danger"))
	message = strings.ReplaceAll(message, "{{veryBad}}", i18n.Get(lang, "veryBad"))
	message = strings.ReplaceAll(message, "{{bad}}", i18n.Get(lang, "bad"))
	message = strings.ReplaceAll(message, "{{normal}}", i18n.Get(lang, "normal"))
	message = strings.ReplaceAll(message, "{{good}}", i18n.Get(lang, "good"))
	message = strings.ReplaceAll(message, "{{dataError}}", i18n.Get(lang, "dataError"))
	message = strings.ReplaceAll(message, "{{dangerAlertLabel}}", i18n.Get(lang, "dangerAlertLabel"))
	message = strings.ReplaceAll(message, "{{dangerAlertEnd}}", i18n.Get(lang, "dangerAlertEnd"))
	_, err = bot.PushMessage(userId, linebot.NewTextMessage(message)).Do()
	if err != nil {
		log.Println(err)
		return
	}
}
