package notificationserviceserver

import (
	"context"
	"encoding/json"
	"errors"

	delivernotifmodel "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/models/deliver_notif_model"
	resetnotifmodel "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/models/reset_notif_model"
	sellnotifmodel "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/models/sell_notif_model"
	pb "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/services/notification_service"
	rabbitmq "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/utils/RabbitMQ"
	"github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnimplementedNotificationServiceServer
	Rmq *rabbitmq.RabbitMQ
}

func (s *Server) ResetNotif(ctx context.Context, req *pb.ResetRequest) (*pb.Response, error) {
	var resetNotif resetnotifmodel.ResetNotification
	var errCode int
	resetNotif.Code = int(req.ResetCode)
	resetNotif.Email = req.Email

	errCode = resetNotif.GetTemplate()
	if errCode != 200 {
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error get reset template",
		}, errors.New("error get reset template")
	}

	body, err := json.Marshal(resetNotif)
	if err != nil {
		logrus.Error("Error marshalling notification: ", err)
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error marshalling reset notification",
		}, errors.New("error marshalling reset notification")
	}

	err = s.Rmq.Publish(body, "reset_email")
	if err != nil {
		logrus.Error("Failed to publish message to RabbitMQ: ", err)
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error send reset email",
		}, errors.New("error send reset email")
	}

	return &pb.Response{
		Code:    int32(200),
		Message: "Success send reset email",
	}, nil
}

func (s *Server) DeliverNotif(ctx context.Context, req *pb.DeliverRequest) (*pb.Response, error) {
	var deliverNotif delivernotifmodel.DeliverNotification
	var errCode int

	deliverNotif.Email = req.Email
	deliverNotif.Title = req.Product
	deliverNotif.Item = req.Item

	errCode = deliverNotif.GetTemplate()
	if errCode != 200 {
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error get deliver template",
		}, errors.New("error get deliver template")
	}

	body, err := json.Marshal(deliverNotif)
	if err != nil {
		logrus.Error("Error marshalling deliver notification: ", err)
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error marshalling deliver notification",
		}, errors.New("error marshalling deliver notification")
	}

	err = s.Rmq.Publish(body, "deliver_email")
	if err != nil {
		logrus.Error("Failed to publish message to RabbitMQ: ", err)
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error send deliver email",
		}, errors.New("error send deliver email")
	}

	return &pb.Response{
		Code:    int32(200),
		Message: "Success send deliver email",
	}, nil
}

func (s *Server) SellNotif(ctx context.Context, req *pb.SellRequest) (*pb.Response, error) {
	var sellNotif sellnotifmodel.SellNotification
	var errCode int

	sellNotif.Email = req.Email
	sellNotif.Title = req.Product
	sellNotif.Price = req.Price
	sellNotif.Fee = req.Fee

	errCode = sellNotif.GetTemplate()
	if errCode != 200 {
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error get sell notif template",
		}, errors.New("error get sell notif template")
	}

	body, err := json.Marshal(sellNotif)
	if err != nil {
		logrus.Error("Error marshalling sell notification: ", err)
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error marshalling sell notification",
		}, errors.New("error marshalling sell notification")
	}

	err = s.Rmq.Publish(body, "sell_email")
	if err != nil {
		logrus.Error("Failed to publish message to RabbitMQ: ", err)
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error send sell email",
		}, errors.New("error send sell email")
	}

	return &pb.Response{
		Code:    int32(200),
		Message: "Success send sell notif email",
	}, nil
}
