package util

import (
	"allto/config"
	"fmt"
	"net/smtp"
)
func SendVerifyEmail(to,verifyLink string) error{
	//读config配置
	host:=config.AppConfig.SMTP.Host
	port:=config.AppConfig.Server.Port
	username:=config.AppConfig.SMTP.Username
	password:=config.AppConfig.SMTP.Password
	from:=config.AppConfig.SMTP.From
	//写邮件
	subject := "AllTo 博客 - 邮箱验证"
      body := fmt.Sprintf("点击下面的链接完成邮箱验证：\n\n%s\n\n如果这不是你的操作，请忽略此邮件。", verifyLink)
      msg := "From: " + from + "\r\n" +
          "To: " + to + "\r\n" +
          "Subject: " + subject + "\r\n" +
          "Content-Type: text/plain; charset=UTF-8\r\n" +
          "\r\n" +
          body

      addr := fmt.Sprintf("%s:%s", host, port)
      auth := smtp.PlainAuth("", username, password, host)

      // 587 端口需要 STARTTLS
      return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}