package Line

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"os"
)

var (
	bot               *linebot.Client
	notificationRange = [6]int{1, 2, 4, 8, 12, 24}
	s                 setting
)

type setting struct {
	LineChannelSecret      string `json:"lineChannelSecret"`
	LineChannelAccessToken string `json:"lineChannelAccessToken"`
}

func init() {
	if file, err := os.ReadFile("setting.json"); err != nil {
		panic(err)
	} else {
		err = json.Unmarshal(file, &s)
		if err != nil {
			panic(err)
		}
	}

	if err := func() (err error) {
		bot, err = linebot.New(s.LineChannelSecret, s.LineChannelAccessToken)
		return
	}(); err != nil {
		panic(err)
	}
}
