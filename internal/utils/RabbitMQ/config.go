package rabbitmq

import (
	"encoding/json"

	notificationmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/notification_model"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	Consumer   <-chan amqp.Delivery
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

func (rmq *RabbitMQ) InitConsumer(queueName string) error {
	if err := rmq.DeclareQueue(queueName); err != nil {
		logrus.Errorln("Error declaring queue: ", err)
		return err
	}

	var err error
	rmq.Consumer, err = rmq.channel.Consume(
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
	return err
}

func (rmq *RabbitMQ) ConsumeNotifDish() {
	for d := range rmq.Consumer {
		var notif notificationmodel.Notification
		err := json.Unmarshal(d.Body, &notif)
		if err != nil {
			logrus.Errorln("Error decoding notification: ", err)
			continue
		}
		notif.SendDish()
	}
}

func (rmq *RabbitMQ) ConsumeNotifMessage() {
	for d := range rmq.Consumer {
		var notif notificationmodel.Notification
		err := json.Unmarshal(d.Body, &notif)
		if err != nil {
			logrus.Errorln("Error decoding notification: ", err)
			continue
		}
		notif.SendMessage()
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
