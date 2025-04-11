package dishmodel

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Dish struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Type     Type     `json:"type"`
	Category Category `json:"category"`
	Tag      Tag      `json:"tag"`
}

type Category struct {
	Id           int
	CategoryDish string `json:"category_dish"`
}

type Tag struct {
	Id      int
	TagDish string `json:"tag_dish"`
}

type Type struct {
	Id       int
	TypeDish string `json:"type_dish"`
}

func (dish *Dish) DecodeFromContext(c *gin.Context) error {
	if err := c.ShouldBindJSON(&dish); err != nil {
		logrus.Error("Error decode JSON: ", err)
		return err
	}
	return nil
}
