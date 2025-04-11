package rabbitmq

import (
	"encoding/json"
	"fmt"

	delivernotifmodel "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/models/deliver_notif_model"
	resetnotifmodel "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/models/reset_notif_model"
	sellnotifmodel "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/models/sell_notif_model"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	connection      *amqp.Connection
	channel         *amqp.Channel
	ResetConsumer   <-chan amqp.Delivery
	DeliverConsumer <-chan amqp.Delivery
	SellConsumer    <-chan amqp.Delivery
}

func (rmq *RabbitMQ) InitConnection() error {
	var err error
	//dsn := os.Getenv("RABBITMQ_URL")
	dsn := "amqp://guest:guest@localhost:5672/"
	rmq.connection, err = amqp.Dial(dsn)
	if err != nil {
		logrus.Errorln("RabbitMQ.connecction: ", err)
		return err
	}
	logrus.Infoln("RabbitMQ.chanel SUCCESS CONNECT")
	return nil
}

func (rmq *RabbitMQ) InitChannel() error {
	var err error
	rmq.channel, err = rmq.connection.Channel()
	if err != nil {
		logrus.Errorln("RabbitMQ.chanel: ", err)
		return err
	}
	logrus.Infoln("RabbitMQ.chanel SUCCESS OPEN")
	return nil
}

func (rmq *RabbitMQ) DeclareQueue(queueName string) error {
	_, err := rmq.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (rmq *RabbitMQ) InitConsumer(queueName string, consumerType string) error {
	if err := rmq.DeclareQueue(queueName); err != nil {
		logrus.Errorln("Error declaring queue: ", err)
		return err
	}

	var err error
	consumer, err := rmq.channel.Consume(
		queueName, // queue name
		"",        // consumer tag
		true,      // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		logrus.Errorln("RabbitMQ.Queue: ", err)
		return err
	}

	switch consumerType {
	case "reset":
		rmq.ResetConsumer = consumer
	case "deliver":
		rmq.DeliverConsumer = consumer
	case "sell":
		rmq.SellConsumer = consumer
	default:
		return fmt.Errorf("unknown consumer type: %s", consumerType)
	}

	return nil
}

func (rmq *RabbitMQ) Publish(body []byte, queueName string) error {
	err := rmq.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		logrus.Errorln("RabbitMQ.Publish: ", err)
	}
	logrus.Infoln("RabbitMQ.Publish: Success publish")
	return err
}

func (rmq *RabbitMQ) ConsumeReset() {
	for d := range rmq.ResetConsumer {
		var notif resetnotifmodel.ResetNotification
		err := json.Unmarshal(d.Body, &notif)
		if err != nil {
			logrus.Errorln("Error decoding notification: ", err)
			continue
		}
		notif.Send()
		logrus.Infoln("RabbitMQ.ConsumeReset: Success send email")
	}
}

func (rmq *RabbitMQ) ConsumeDeliver() {
	for d := range rmq.DeliverConsumer {
		var notif delivernotifmodel.DeliverNotification
		err := json.Unmarshal(d.Body, &notif)
		if err != nil {
			logrus.Errorln("Error decoding notification: ", err)
			continue
		}
		notif.Send()
	}
}

func (rmq *RabbitMQ) ConsumeSell() {
	for d := range rmq.SellConsumer {
		var notif sellnotifmodel.SellNotification
		err := json.Unmarshal(d.Body, &notif)
		if err != nil {
			logrus.Errorln("Error decoding notification: ", err)
			continue
		}
		notif.Send()
	}
}

func (rmq *RabbitMQ) Close() {
	if rmq.channel != nil {
		rmq.channel.Close()
	}
	if rmq.connection != nil {
		rmq.connection.Close()
	}
}
