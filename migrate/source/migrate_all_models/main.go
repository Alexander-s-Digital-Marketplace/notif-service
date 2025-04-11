package main

import (
	loggerconfig "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/config/logger"
	"github.com/Alexander-s-Digital-Marketplace/notif-service/internal/database"
	templatemodel "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/models/template_model"
)

func main() {
	loggerconfig.Init()

	var db database.DataBase
	db.InitDB()

	var template templatemodel.Template
	template.MigrateToDB(db)

	db.CloseDB()
}
