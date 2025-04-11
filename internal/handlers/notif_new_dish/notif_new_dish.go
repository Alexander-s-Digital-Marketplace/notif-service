package notifnewdish

import (
	"encoding/json"

	dishmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/dish_model"
	emailmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/email_model"
	notificationmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/notification_model"
	rabbitmq "github.com/KusakinDev/Catering-Notif-Service/internal/utils/RabbitMQ"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NotifNewDish(c *gin.Context, rmq *rabbitmq.RabbitMQ) (int, string) {
	var dish dishmodel.Dish
	dish.DecodeFromContext(c)

	var notif notificationmodel.Notification
	notif.GetTemplateByTag("new_dish")
	notif.Dish = dish

	var email emailmodel.Email
	emails, code := email.GetAllFromTable()
	if code != 200 {
		return code, "Error get all emails from table"
	}

	for _, email := range emails {
		notif.Email = email

		body, err := json.Marshal(notif)
		if err != nil {
			logrus.Error("Error marshalling notification: ", err)
			return 500, "Internal error"
		}

		err = rmq.Publish(body, "dishQueue")
		if err != nil {
			logrus.Error("Failed to publish message to RabbitMQ: ", err)
			return 500, "Error sending notification"
		}
	}
	return 200, "Success broadcast"
}
