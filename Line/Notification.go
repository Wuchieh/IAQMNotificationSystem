package Line

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Wuchieh/IAQMNotificationSystem/Database"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	salt         = "511b57761616b978a02fb4f4a90b8d05"
	expectedHash = "039f8aaac8ef9ac536cba9dd5e584d5854e11b9325fae0a518ef3cb4c7675de4"
)

func SendPushNotification(c *gin.Context) {
	token := c.PostForm("token") // POST 中拿 token
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token"})
		return
	}

	// 對token加鹽
	hashedToken := hashWithSalt(token, salt)
	fmt.Println(hashedToken)
	// 驗證加鹽後的結果是否與預期相符
	if hashedToken != expectedHash {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	lineID := c.PostForm("lineID") // 從 POST 中取 lineID
	if lineID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing lineID"})
		return
	}

	// 檢查 lineID 是否在資料庫中
	user, err := Database.GetUserByLineID(lineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	msg := c.PostForm("message") // 從POST取得訊息
	if msg == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message can not empty"})
		return
	}
	SendMessage(user.LineId, msg)

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent"})
}

func hashWithSalt(value string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(value + salt))
	hashedValue := hash.Sum(nil)
	return hex.EncodeToString(hashedValue)
}
