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

	sqlStatements := []string{
		`INSERT INTO emails (first_name, last_name, email, access) VALUES
        ('Alex', 'Kus', 'AlexKus@ex.com', 1),
        ('Alice', 'Deeb', 'AliceDeeb@ex.com', 1),
        ('Bob', 'Qwe', 'BobQwe@ex.com', 1);`,
	}

	for _, stmt := range sqlStatements {
		if err := db.Connection.Exec(stmt).Error; err != nil {
			logrus.Println("Error executing seed: ", stmt, err)
			return
		}
	}

	logrus.Println("Success seeding")

	db.CloseDB()
}
