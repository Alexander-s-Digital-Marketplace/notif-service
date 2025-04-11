package grpcnotifnewmenu

import (
	"context"
	"encoding/json"
	"errors"

	emailmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/email_model"
	notificationmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/notification_model"
	pb "github.com/KusakinDev/Catering-Notif-Service/internal/services/notif_new_menu/notif_new_menu"
	rabbitmq "github.com/KusakinDev/Catering-Notif-Service/internal/utils/RabbitMQ"
	"github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnimplementedNotifNewMenuServiceServer
	Rmq *rabbitmq.RabbitMQ
}

func (s *Server) NotifNewMenu(ctx context.Context, req *pb.MenuRequest) (*pb.MenuResponse, error) {
	message := req.Message

	var email emailmodel.Email
	emails, code := email.GetAllFromTable()
	if code != 200 {
		return &pb.MenuResponse{Message: "Error get all emails from table"}, errors.New("error get all emails from table")
	}

	var notif notificationmodel.Notification
	notif.GetTemplateByTag("new_menu")
	notif.Message = message
	for _, email := range emails {
		notif.Email = email

		body, err := json.Marshal(notif)
		if err != nil {
			logrus.Error("Error marshalling notification: ", err)
			return &pb.MenuResponse{Message: "Error marshalling notification"}, errors.New("error marshalling notification")
		}

		err = s.Rmq.Publish(body, "menuQueue")
		if err != nil {
			logrus.Error("Failed to publish message to RabbitMQ: ", err)
			return &pb.MenuResponse{Message: "Failed to publish message to RabbitMQ"}, errors.New("failed to publish message to rabbit_mq")
		}
	}

	return &pb.MenuResponse{Message: "Success broadcast"}, nil
}
