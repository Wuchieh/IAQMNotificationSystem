package Server

import (
	"github.com/Wuchieh/IAQMNotificationSystem/Line"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.POST("/danger", dangerAlerts)
	v1.POST("/send-notification", Line.SendPushNotification) // 新增路由處理推播訊息的請求
}