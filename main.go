package main

import (
	"database/sql"
	"github.com/Wuchieh/IAQMNotificationSystem/Database"
	"github.com/Wuchieh/IAQMNotificationSystem/Line"
	"github.com/Wuchieh/IAQMNotificationSystem/Server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	Database.DatabaseInit()
	if a := <-Database.Sign; a == 0 {
		log.Println("資料庫連接成功")
		defer func(Db *sql.DB) {
			err := Db.Close()
			if err != nil {
				log.Println("資料必關閉連線異常:", err)
			} else {
				log.Println("資料庫已中斷連線")
			}
		}(Database.Db)
	} else {
		log.Println("資料庫連接失敗")
	}
	sc := make(chan os.Signal, 1)
	go Line.NotificationCronjob()
	go func() {
		Server.Run()
		sc <- syscall.SIGINT
	}()

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
