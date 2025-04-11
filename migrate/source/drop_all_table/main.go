package main

import (
	loggerconfig "github.com/KusakinDev/Catering-Notif-Service/internal/config/logger"
	"github.com/KusakinDev/Catering-Notif-Service/internal/database"
	"github.com/sirupsen/logrus"
)

func main() {
	loggerconfig.Init()

	var db database.DataBase
	db.InitDB()

	query := []string{
		`DROP TABLE emails CASCADE;`,
	}

	for _, stmt := range query {
		if err := db.Connection.Exec(stmt).Error; err != nil {
			logrus.Println("Error executing drop: ", stmt, err)
		}
	}

	logrus.Println("All table is droped")

	db.CloseDB()
}
