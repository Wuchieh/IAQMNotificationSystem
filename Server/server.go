package Server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

var (
	s setting
)

type setting struct {
	ServerAddr string `json:"serverAddr"`
	RUNMODE    string `json:"RUN MODE"`
}

func init() {
	file, err := os.ReadFile("setting.json")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func Run() {
	setGinMode(s.RUNMODE)
	r := gin.Default()
	Router(r)

	log.Println(r.Run(s.ServerAddr))
}

func setGinMode(mode string) {
	if len(mode) < 1 {
		return
	}
	m := strings.ToLower(mode[:1])
	switch m {
	case "r":
		gin.SetMode(gin.ReleaseMode)
	case "t":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
