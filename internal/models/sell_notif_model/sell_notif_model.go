package sellnotifmodel

import (
	"bytes"
	"html/template"

	emailconfig "github.com/Alexander-s-Digital-Marketplace/notif-service/.env/email"
	templatemodel "github.com/Alexander-s-Digital-Marketplace/notif-service/internal/models/template_model"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type SellNotification struct {
	Email    string
	Title    string
	Price    float64
	Fee      float64
	Template templatemodel.Template
}

func (notif *SellNotification) GetTemplate() int {
	notif.Template.Description = "notif_sell"
	code := notif.Template.GetFromTableByDescription()
	if code != 200 {
		logrus.Errorln("Error find template")
		return 404
	}
	return 200
}

func (notif *SellNotification) Send() int {

	tmpl, err := template.New("notif_sell").Parse(notif.Template.Template)
	if err != nil {
		logrus.Errorln(err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, notif)
	if err != nil {
		logrus.Errorln(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailconfig.Email)
	m.SetHeader("To", notif.Email)
	m.SetHeader("Subject", "ADM: Successful sale!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(emailconfig.Host, emailconfig.Port, emailconfig.Email, emailconfig.Password)

	err = d.DialAndSend(m)
	if err != nil {
		logrus.Error("Error send email: ", err)
		return 404
	}
	return 200
}
