package gomail

// import (
// 	"bytes"
// 	"crypto/tls"
// 	"fmt"
// 	"hospitality-service/domain"
// 	"html/template"
// 	"log/slog"

// 	"github.com/spf13/viper"
// 	"gopkg.in/gomail.v2"
// )

// type gomailRepository struct{}

// func NewGoMailRepository() domain.SMTPRepository {
// 	return &gomailRepository{}
// }

// func (u *gomailRepository) SendEmailToUser(req domain.SMTPRequest) (status int, err error) {
// 	slog.Debug("[Repository][SendEmail]")

// 	var bodyHTML bytes.Buffer
// 	ht, err := template.ParseFiles("./html_template/" + req.TemplateName)
// 	if err != nil {
// 		slog.Error("[Repository][SendEmail] Err", ":", err.Error())
// 		status = domain.StatusInternalServerError
// 		return
// 	}
// 	err = ht.Execute(&bodyHTML, req.Parameter)
// 	if err != nil {
// 		slog.Error("[Repository][SendEmail] Err", ":", err.Error())
// 		status = domain.StatusInternalServerError
// 		return
// 	}

// 	message := gomail.NewMessage()
// 	message.SetHeader("From", viper.GetString("smtp.email_sender"))
// 	message.SetHeader("To", req.Recipient)
// 	message.SetHeader("Subject", req.Subject)

// 	for contentID, imagePath := range req.EmbedImage {
// 		message.Embed(imagePath, gomail.SetHeader(map[string][]string{
// 			"Content-ID": {fmt.Sprintf("<%s>", contentID)},
// 		}))
// 	}

// 	message.SetBody("text/html", bodyHTML.String())

// 	d := gomail.NewDialer(viper.GetString("smtp.mail"), int(viper.GetInt64("smtp.port")), viper.GetString("smtp.email_sender"), viper.GetString("smtp.password"))

// 	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

// 	tlsConfig := &tls.Config{InsecureSkipVerify: true}
// 	if viper.GetInt64("smtp.start_tls") == 1 {
// 		d.TLSConfig = tlsConfig
// 	} else if viper.GetInt64("smtp.tsl_or_ssl") == 1 {
// 		d.SSL = true
// 	}

// 	if err = d.DialAndSend(message); err != nil {
// 		slog.Error("[Repository][SendEmail] Err", ":", err.Error())
// 		status = domain.StatusInternalServerError
// 		return
// 	}
// 	return
// }
