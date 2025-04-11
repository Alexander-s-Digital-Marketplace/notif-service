package main

import (
	loggerconfig "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/config/logger"
	"github.com/Alexander-s-Digital-Marketplace/notif-service/internal/database"
	"github.com/sirupsen/logrus"
)

func main() {
	loggerconfig.Init()

	var db database.DataBase
	db.InitDB()

	query := []string{
		`DELETE FROM emails;`,
	}

	for _, stmt := range query {
		if err := db.Connection.Exec(stmt).Error; err != nil {
			logrus.Errorln("Error executing clear: ", stmt, err)
		}
	}

	logrus.Infoln("All table is clear")

	db.CloseDB()
}
