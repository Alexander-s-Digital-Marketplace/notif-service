package templatemodel

import (
	"github.com/Alexander-s-Digital-Marketplace/notif-service/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Template struct {
	Id          int    `gorm:"primaryKey;autoIncrement"`
	Template    string `gorm:"type:text"`
	Description string `gorm:"type:varchar(50)"`
}

func (template *Template) DecodeFromContext(c *gin.Context) error {
	if err := c.ShouldBindJSON(&template); err != nil {
		logrus.Error("Error decode JSON: ", err)
		return err
	}
	return nil
}

func (template *Template) AddToTable() int {
	var db database.DataBase
	db.InitDB()

	err := db.Connection.Create(&template).Error
	if err != nil {
		db.CloseDB()
		logrus.Error("Error add to table: ", err)
		return 503
	}

	db.CloseDB()
	return 200
}

func (template *Template) GetFromTableById() int {
	var db database.DataBase
	db.InitDB()

	err := db.Connection.First(&template).Error
	if err != nil {
		db.CloseDB()
		return 503
	}

	db.CloseDB()
	return 200
}

func (template *Template) GetAllFromTable() ([]Template, int) {
	var db database.DataBase
	db.InitDB()

	var templates []Template

	err := db.Connection.Find(&templates).Error
	if err != nil {
		db.CloseDB()
		logrus.Errorln("Error get all templates from table")
		return []Template{}, 503
	}
	db.CloseDB()
	return templates, 200
}

func (template *Template) GetFromTableByDescription() int {
	var db database.DataBase
	db.InitDB()

	err := db.Connection.Where("description = ?", template.Description).First(&template).Error
	if err != nil {
		db.CloseDB()
		return 503
	}

	db.CloseDB()
	return 200
}

func (template *Template) MigrateToDB(db database.DataBase) error {
	err := db.Connection.AutoMigrate(&Template{})
	if err != nil {
		logrus.Errorln("Error migrate template model :")
		return err
	}
	logrus.Infoln("Success migrate template model :")
	return nil
}
