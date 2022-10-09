package services

import (
	"Twitta/global"
	"github.com/go-gomail/gomail"
	"go.uber.org/zap"
)

// 发送邮件服务: params: 主题、发送人名字、发送人地址、接收人名字、接受人地址、内容类型、内容
func SendEmail(topic string, name string, addr string, contentType string, content string) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", global.ServerConfig.Email.SendPersonAddr, global.ServerConfig.Email.SendPersonName)
	m.SetHeader("To", m.FormatAddress(addr, name))
	m.SetHeader("Subject", topic)
	m.SetBody(contentType, content)
	d := gomail.NewDialer(global.ServerConfig.Email.Host, global.ServerConfig.Email.Port, global.ServerConfig.Email.SendPersonAddr, global.ServerConfig.Email.Password)
	if err := d.DialAndSend(m); err != nil {
		go zap.S().Errorf(err.Error())
	} else {
		go zap.S().Info("发送邮件成功")
	}
}

func ConsumeEmailMessage() {
	// 检测是否可以获取到数据，如果能，那么消费，否则阻塞
	for {
		select {
		case emailMessage := <-EmailMessageChan:
			SendEmail(emailMessage.Topic, emailMessage.Name, emailMessage.Addr, emailMessage.ContentType, emailMessage.Content)
		default:
		}
	}
}

type EmailMessage struct {
	Topic       string
	Name        string
	Addr        string
	ContentType string
	Content     string
}

var EmailMessageChan = make(chan *EmailMessage, 20000)
