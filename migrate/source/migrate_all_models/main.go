package main

import (
	loggerconfig "github.com/KusakinDev/Catering-Notif-Service/internal/config/logger"
	"github.com/KusakinDev/Catering-Notif-Service/internal/database"
	emailmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/email_model"
	templatemodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/template_model"
)

func main() {
	loggerconfig.Init()

	var db database.DataBase
	db.InitDB()

	var email emailmodel.Email
	email.MigrateToDB(db)

	var template templatemodel.Template
	template.MigrateToDB(db)

	db.CloseDB()
}
