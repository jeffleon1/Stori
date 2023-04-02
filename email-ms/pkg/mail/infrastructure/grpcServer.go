package infra

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	pb "github.com/jeffleon1/email-ms/pkg/mail/domain/proto"
	"github.com/sirupsen/logrus"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailserver struct {
	pb.UnimplementedMailServiceServer
	SMTPClient *mail.SMTPClient
}

func (m *Mailserver) SendMail(ctx context.Context, req *pb.MailRequest) (*pb.MailResponse, error) {
	err := m.PrepareEmail(req)
	if err != nil {
		logrus.Errorf("Failed")
		panic(err)
	}
	return &pb.MailResponse{
		Message: fmt.Sprintf("your mail was sended succesfull to %s", req.To),
	}, nil
}

func (m *Mailserver) PrepareEmail(msg *pb.MailRequest) error {
	htmlBody, err := m.buildHTMLMessage(msg)
	if err != nil {
		logrus.Errorf("Error whith creation mail body : %s", err)
		return err
	}
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	email.SetBody(mail.TextHTML, htmlBody)
	err = email.Send(m.SMTPClient)
	if err != nil {
		logrus.Errorf("Error sending mail : %s", err)
		return err
	}

	logrus.Infof("Sent Email to %s", msg.To)
	return nil
}

func (m *Mailserver) buildHTMLMessage(msg *pb.MailRequest) (string, error) {
	templateToRender := "./pkg/mail/infrastructure/templates/account_resume.html.gohtml"

	templ, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", nil
	}

	var tpl bytes.Buffer

	if err = templ.ExecuteTemplate(&tpl, "body", msg); err != nil {
		return "", nil
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mailserver) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
