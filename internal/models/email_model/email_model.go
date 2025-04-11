package emailmodel

import (
	"github.com/KusakinDev/Catering-Notif-Service/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Email struct {
	Id        int    `gorm:"primaryKey;autoIncrement"`
	FirstName string `gorm:"type:varchar(50)"`
	LastName  string `gorm:"type:varchar(50)"`
	Email     string `gorm:"type:varchar(50)"`
	Access    int    `gorm:"type:integer"`
}

func (email *Email) DecodeFromContext(c *gin.Context) error {
	if err := c.ShouldBindJSON(&email); err != nil {
		logrus.Error("Error decode JSON: ", err)
		return err
	}
	return nil
}

func (email *Email) AddToTable() int {
	var db database.DataBase
	db.InitDB()

	err := db.Connection.Create(&email).Error
	if err != nil {
		db.CloseDB()
		logrus.Error("Error add to table: ", err)
		return 503
	}

	db.CloseDB()
	return 200
}

func (email *Email) GetFromTableById() int {
	var db database.DataBase
	db.InitDB()

	err := db.Connection.First(&email)
	if err != nil {
		db.CloseDB()
		return 503
	}

	db.CloseDB()
	return 200
}

func (email *Email) GetAllFromTable() ([]Email, int) {
	var db database.DataBase
	db.InitDB()

	var emails []Email

	err := db.Connection.Find(&emails).Error
	if err != nil {
		db.CloseDB()
		logrus.Errorln("Error get all emails from table")
		return []Email{}, 503
	}
	db.CloseDB()
	return emails, 200
}

func (email *Email) GetFromTableByEmail() int {
	var db database.DataBase
	db.InitDB()

	err := db.Connection.Where("email = ?", email.Email).First(&email)
	if err != nil {
		db.CloseDB()
		return 503
	}

	db.CloseDB()
	return 200
}

func (email *Email) MigrateToDB(db database.DataBase) error {
	err := db.Connection.AutoMigrate(&Email{})
	if err != nil {
		logrus.Errorln("Error migrate email model :")
		return err
	}
	logrus.Infoln("Success migrate email model :")
	return nil
}
