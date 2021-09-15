//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package common

import (
	"crypto/tls"
	"github.com/beego/beego/v2/server/web"
	"github.com/xhit/go-simple-mail/v2"
	"log"
	"mogu-go-v2/models"
	"strconv"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/26 10:06 上午
 * @version 1.0
 */

//var mail *utils.Email

/*type config struct {
	Identity string `json:"identity"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	From     string `json:"from"`
}*/

/*func initialize() {
var systemConfig models.SystemConfig
DB.Where("status=?", 1).First(&systemConfig)
username := systemConfig.EmailUserName
password := systemConfig.EmailPassword
host := systemConfig.SmtpAddress
port, _ := strconv.Atoi(systemConfig.SmtpPort)
from := systemConfig.Email*/
/*configJson := config{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		From:     from,
	}
	b, _ := json.Marshal(configJson)
	mail = utils.NewEMail(string(b))
}*/

func sendEmail(email, text string) {
	//initialize()
	var systemConfig models.SystemConfig
	DB.Where("status=?", 1).First(&systemConfig)
	username := systemConfig.EmailUserName
	password := systemConfig.EmailPassword
	host := systemConfig.SmtpAddress
	port, _ := strconv.Atoi(systemConfig.SmtpPort)
	from := systemConfig.Email
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = host
	server.Port = port
	server.Username = username
	server.Password = password
	server.Encryption = mail.EncryptionSSL

	// Since v2.3.0 you can specified authentication type:
	// - PLAIN (default)
	// - LOGIN
	// - CRAM-MD5
	// server.Authentication = mail.AuthPlain

	// Variable to keep alive connection
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}
	/*mail.Subject = "蘑菇博客"
	mail.To = append([]string{}, email)
	mail.HTML = text
	err := mail.Send()
	if err != nil {
		panic(err)
	}*/
	eMail := mail.NewMSG()
	eMail.SetFrom(from).
		AddTo(email).
		SetSubject("蘑菇博客GO")

	eMail.SetBody(mail.TextHTML, text)

	// Call Send and pass the client
	err = eMail.Send(smtpClient)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent")
	}
}

type email struct{}

var projectName, _ = web.AppConfig.String("project_name")
var projectNameEn, _ = web.AppConfig.String("project_name_en")
var logo, _ = web.AppConfig.String("logo")
var dataWebsiteUrl, _ = web.AppConfig.String("data_website_url")
var dataWebUrl, _ = web.AppConfig.String("data_web_url")

func (*email) SendRegisterEmail(user models.User, token string) {
	text := "<html>\r\n" +
		" <head>\r\n" +
		"  <title>" + projectName + "</title>\r\n" +
		" </head>\r\n" +
		" <body>\r\n" +
		"  <div id=\"contentDiv\" onmouseover=\"getTop().stopPropagation(event);\" onclick=\"getTop().preSwapLink(event, 'spam', 'ZC1222-PrLAp4T0Z7Z7UUMYzqLkb8a');\" style=\"position:relative;font-size:14px;height:auto;padding:15px 15px 10px 15px;z-index:1;zoom:1;line-height:1.7;\" class=\"body\">    \r\n" +
		"  <div id=\"qm_con_body\"><div id=\"mailContentContainer\" class=\"qmbox qm_con_body_content qqmail_webmail_only\" style=\"\">\r\n" +
		"<style>\r\n" +
		"  .qmbox .email-body{color:#40485B;font-size:14px;font-family:-apple-system, \"Helvetica Neue\", Helvetica, \"Nimbus Sans L\", \"Segoe UI\", Arial, \"Liberation Sans\", \"PingFang SC\", \"Microsoft YaHei\", \"Hiragino Sans GB\", \"Wenquanyi Micro Hei\", \"WenQuanYi Zen Hei\", \"ST Heiti\", SimHei, \"WenQuanYi Zen Hei Sharp\", sans-serif;background:#f8f8f8;}.qmbox .pull-right{float:right;}.qmbox a{color:#FE7300;text-decoration:underline;}.qmbox a:hover{color:#fe9d4c;}.qmbox a:active{color:#b15000;}.qmbox .logo{text-align:center;margin-bottom:20px;}.qmbox .panel{background:#fff;border:1px solid #E3E9ED;margin-bottom:10px;}.qmbox .panel-header{font-size:18px;line-height:30px;padding:10px 20px;background:#fcfcfc;border-bottom:1px solid #E3E9ED;}.qmbox .panel-body{padding:20px;}.qmbox .container{width:50%;min-width:600px;padding:20px;margin:0 auto;}.qmbox .text-center{text-align:center;}.qmbox .thumbnail{padding:4px;max-width:100%;border:1px solid #E3E9ED;}.qmbox .btn-primary{color:#fff;font-size:16px;padding:8px 14px;line-height:20px;border-radius:2px;display:inline-block;background:#FE7300;text-decoration:none;}.qmbox .btn-primary:hover,.qmbox .btn-primary:active{color:#fff;}.qmbox .footer{color:#9B9B9B;font-size:12px;margin-top:40px;}.qmbox .footer a{color:#9B9B9B;}.qmbox .footer a:hover{color:#fe9d4c;}.qmbox .footer a:active{color:#b15000;}.qmbox .email-body#mail_to_teacher{line-height:26px;color:#40485B;font-size:16px;padding:0px;}.qmbox .email-body#mail_to_teacher .container,.qmbox .email-body#mail_to_teacher .panel-body{padding:0px;}.qmbox .email-body#mail_to_teacher .container{padding-top:20px;}.qmbox .email-body#mail_to_teacher .textarea{padding:32px;}.qmbox .email-body#mail_to_teacher .say-hi{font-weight:500;}.qmbox .email-body#mail_to_teacher .paragraph{margin-top:24px;}.qmbox .email-body#mail_to_teacher .paragraph .pro-name{color:#000000;}.qmbox .email-body#mail_to_teacher .paragraph.link{margin-top:32px;text-align:center;}.qmbox .email-body#mail_to_teacher .paragraph.link .button{background:#4A90E2;border-radius:2px;color:#FFFFFF;text-decoration:none;padding:11px 17px;line-height:14px;display:inline-block;}.qmbox .email-body#mail_to_teacher ul.pro-desc{list-style-type:none;margin:0px;padding:0px;padding-left:16px;}.qmbox .email-body#mail_to_teacher ul.pro-desc li{position:relative;}.qmbox .email-body#mail_to_teacher ul.pro-desc li::before{content:'';width:3px;height:3px;border-radius:50%;background:red;position:absolute;left:-15px;top:11px;background:#40485B;}.qmbox .email-body#mail_to_teacher .blackboard-area{height:600px;padding:40px;background-image:url();color:#FFFFFF;}.qmbox .email-body#mail_to_teacher .blackboard-area .big-title{font-size:32px;line-height:45px;text-align:center;}.qmbox .email-body#mail_to_teacher .blackboard-area .desc{margin-top:8px;}.qmbox .email-body#mail_to_teacher .blackboard-area .desc p{margin:0px;text-align:center;line-height:28px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card:nth-child(odd){float:left;margin-top:45px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card:nth-child(even){float:right;margin-top:45px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card .title{font-size:18px;text-align:center;margin-bottom:10px;}\r\n" +
		"</style>\r\n" +
		"<meta>\r\n" +
		"<div class=\"email-body\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
		"<div class=\"container\">\r\n" +
		"<div class=\"logo\">\r\n" +
		"<img src=\"" + logo + "\",height=\"100\" width=\"100\">\r\n" +
		"</div>\r\n" +
		"<div class=\"panel\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
		"<div class=\"panel-header\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
		projectName + "邮箱绑定\r\n" +
		"\r\n" +
		"</div>\r\n" +
		"<div class=\"panel-body\">\r\n" +
		"<p>您好 <a href=\"mailto:" + user.Email + "\" rel=\"noopener\" target=\"_blank\">" + user.NickName + "<wbr></a>！</p>\r\n" +
		"<p>欢迎您给" + projectName + "账号绑定邮箱，请点击下方链接进行绑定</p>\r\n" +
		"<p>地址：" + "<a href=\"" + dataWebUrl + "/oauth/bindUserEmail/" + token + "/" + user.ValidCode + "\">点击这里</a>" + "</p>\r\n" +
		"\r\n" +
		"</div>\r\n" +
		"</div>\r\n" +
		"<div class=\"footer\">\r\n" +
		"<a href=\" " + dataWebsiteUrl + "\">@" + projectNameEn + "</a>\n" +
		"<div class=\"pull-right\"></div>\r\n" +
		"</div>\r\n" +
		"</div>\r\n" +
		"</div>\r\n" +
		"<style type=\"text/css\">.qmbox style, .qmbox script, .qmbox head, .qmbox link, .qmbox meta {display: none !important;}</style></div></div><!-- --><style>#mailContentContainer .txt {height:auto;}</style>  </div>\r\n" +
		" </body>\r\n" +
		"</html>"
	sendEmail(user.Email, text)
}

func (*email) SentCommentEmail(m map[string]string) {
	email := m["email"]
	text := m["text"]
	toText := m["to_text"]
	nickName := m["nickname"]
	toUserNickName := m["to_nickname"]
	userUid := m["user_uid"]
	url := m["url"]
	content :=
		"<html>\r\n" +
			" <head>\r\n" +
			"  <title> " + projectName + "</title>\r\n" +
			" </head>\r\n" +
			" <body>\r\n" +
			"  <div id=\"contentDiv\" onmouseover=\"getTop().stopPropagation(event);\" onclick=\"getTop().preSwapLink(event, 'spam', 'ZC1222-PrLAp4T0Z7Z7UUMYzqLkb8a');\" style=\"position:relative;font-size:14px;height:auto;padding:15px 15px 10px 15px;z-index:1;zoom:1;line-height:1.7;\" class=\"body\">    \r\n" +
			"  <div id=\"qm_con_body\"><div id=\"mailContentContainer\" class=\"qmbox qm_con_body_content qqmail_webmail_only\" style=\"\">\r\n" +
			"<style>\r\n" +
			"  .qmbox .email-body{color:#40485B;font-size:14px;font-family:-apple-system, \"Helvetica Neue\", Helvetica, \"Nimbus Sans L\", \"Segoe UI\", Arial, \"Liberation Sans\", \"PingFang SC\", \"Microsoft YaHei\", \"Hiragino Sans GB\", \"Wenquanyi Micro Hei\", \"WenQuanYi Zen Hei\", \"ST Heiti\", SimHei, \"WenQuanYi Zen Hei Sharp\", sans-serif;background:#f8f8f8;}.qmbox .pull-right{float:right;}.qmbox a{color:#FE7300;text-decoration:underline;}.qmbox a:hover{color:#fe9d4c;}.qmbox a:active{color:#b15000;}.qmbox .logo{text-align:center;margin-bottom:20px;}.qmbox .panel{background:#fff;border:1px solid #E3E9ED;margin-bottom:10px;}.qmbox .panel-header{font-size:18px;line-height:30px;padding:10px 20px;background:#fcfcfc;border-bottom:1px solid #E3E9ED;}.qmbox .panel-body{padding:20px;}.qmbox .container{width:50%;min-width:300px;padding:20px;margin:0 auto;}.qmbox .text-center{text-align:center;}.qmbox .thumbnail{padding:4px;max-width:100%;border:1px solid #E3E9ED;}.qmbox .btn-primary{color:#fff;font-size:16px;padding:8px 14px;line-height:20px;border-radius:2px;display:inline-block;background:#FE7300;text-decoration:none;}.qmbox .btn-primary:hover,.qmbox .btn-primary:active{color:#fff;}.qmbox .footer{color:#9B9B9B;font-size:12px;margin-top:40px;}.qmbox .footer a{color:#9B9B9B;}.qmbox .footer a:hover{color:#fe9d4c;}.qmbox .footer a:active{color:#b15000;}.qmbox .email-body#mail_to_teacher{line-height:26px;color:#40485B;font-size:16px;padding:0px;}.qmbox .email-body#mail_to_teacher .container,.qmbox .email-body#mail_to_teacher .panel-body{padding:0px;}.qmbox .email-body#mail_to_teacher .container{padding-top:20px;}.qmbox .email-body#mail_to_teacher .textarea{padding:32px;}.qmbox .email-body#mail_to_teacher .say-hi{font-weight:500;}.qmbox .email-body#mail_to_teacher .paragraph{margin-top:24px;}.qmbox .email-body#mail_to_teacher .paragraph .pro-name{color:#000000;}.qmbox .email-body#mail_to_teacher .paragraph.link{margin-top:32px;text-align:center;}.qmbox .email-body#mail_to_teacher .paragraph.link .button{background:#4A90E2;border-radius:2px;color:#FFFFFF;text-decoration:none;padding:11px 17px;line-height:14px;display:inline-block;}.qmbox .email-body#mail_to_teacher ul.pro-desc{list-style-type:none;margin:0px;padding:0px;padding-left:16px;}.qmbox .email-body#mail_to_teacher ul.pro-desc li{position:relative;}.qmbox .email-body#mail_to_teacher ul.pro-desc li::before{content:'';width:3px;height:3px;border-radius:50%;background:red;position:absolute;left:-15px;top:11px;background:#40485B;}.qmbox .email-body#mail_to_teacher .blackboard-area{height:600px;padding:40px;background-image:url();color:#FFFFFF;}.qmbox .email-body#mail_to_teacher .blackboard-area .big-title{font-size:32px;line-height:45px;text-align:center;}.qmbox .email-body#mail_to_teacher .blackboard-area .desc{margin-top:8px;}.qmbox .email-body#mail_to_teacher .blackboard-area .desc p{margin:0px;text-align:center;line-height:28px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card:nth-child(odd){float:left;margin-top:45px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card:nth-child(even){float:right;margin-top:45px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card .title{font-size:18px;text-align:center;margin-bottom:10px;}\r\n" +
			"</style>\r\n" +
			"<meta>\r\n" +
			"<div class=\"email-body\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
			"<div class=\"container\">\r\n" +
			"<div class=\"logo\">\r\n" +
			"<img src=\"" + logo + "\",height=\"100\" width=\"100\">\r\n" +
			"</div>\r\n" +
			"<div class=\"panel\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
			"<div class=\"panel-header\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
			"评论提醒\r\n" +
			"\r\n" +
			"</div>\r\n" +
			"<div class=\"panel-body\">\r\n" +
			"<p>您好 <a href=\"mailto:" + email + "\" rel=\"noopener\" target=\"_blank\">" + toUserNickName + "<wbr></a>！</p>\r\n" +
			"<p>" + nickName + " 对您的评论：" + "<a href=\"" + url + "\">" + toText + "</a>" + "   进行了回复</p>\r\n" +
			"\r\n" +
			"<p>回复内容为：" + "<a href=\"" + url + "\">" + text + "</a>" + "</p>\r\n" +
			"\r\n" +
			"<p>如果邮件通知干扰了您，可以点击右侧链接关闭通知：" + "<a href=\"" + dataWebUrl + "/web/comment/closeEmailNotification/" + userUid + "\">点击这里</a>" + "</p>\r\n" +
			"</div>\r\n" +
			"</div>\r\n" +
			"<div class=\"footer\">\r\n" +
			"<a href=\" " + dataWebsiteUrl + "\">@" + projectNameEn + "</a>\n" +
			"<div class=\"pull-right\"></div>\r\n" +
			"</div>\r\n" +
			"</div>\r\n" +
			"</div>\r\n" +
			"<style type=\"text/css\">.qmbox style, .qmbox script, .qmbox head, .qmbox link, .qmbox meta {display: none !important;}</style></div></div><!-- --><style>#mailContentContainer .txt {height:auto;}</style>  </div>\r\n" +
			" </body>\r\n" +
			"</html>"
	sendEmail(email, content)
}

func (c *email) SentSimpleEmail(email string, text string) {
	content :=
		"<html>\r\n" +
			" <head>\r\n" +
			"  <title> " + projectName + "</title>\r\n" +
			" </head>\r\n" +
			" <body>\r\n" +
			"  <div id=\"contentDiv\" onmouseover=\"getTop().stopPropagation(event);\" onclick=\"getTop().preSwapLink(event, 'spam', 'ZC1222-PrLAp4T0Z7Z7UUMYzqLkb8a');\" style=\"position:relative;font-size:14px;height:auto;padding:15px 15px 10px 15px;z-index:1;zoom:1;line-height:1.7;\" class=\"body\">    \r\n" +
			"  <div id=\"qm_con_body\"><div id=\"mailContentContainer\" class=\"qmbox qm_con_body_content qqmail_webmail_only\" style=\"\">\r\n" +
			"<style>\r\n" +
			"  .qmbox .email-body{color:#40485B;font-size:14px;font-family:-apple-system, \"Helvetica Neue\", Helvetica, \"Nimbus Sans L\", \"Segoe UI\", Arial, \"Liberation Sans\", \"PingFang SC\", \"Microsoft YaHei\", \"Hiragino Sans GB\", \"Wenquanyi Micro Hei\", \"WenQuanYi Zen Hei\", \"ST Heiti\", SimHei, \"WenQuanYi Zen Hei Sharp\", sans-serif;background:#f8f8f8;}.qmbox .pull-right{float:right;}.qmbox a{color:#FE7300;text-decoration:underline;}.qmbox a:hover{color:#fe9d4c;}.qmbox a:active{color:#b15000;}.qmbox .logo{text-align:center;margin-bottom:20px;}.qmbox .panel{background:#fff;border:1px solid #E3E9ED;margin-bottom:10px;}.qmbox .panel-header{font-size:18px;line-height:30px;padding:10px 20px;background:#fcfcfc;border-bottom:1px solid #E3E9ED;}.qmbox .panel-body{padding:20px;}.qmbox .container{width:50%;min-width:300px;padding:20px;margin:0 auto;}.qmbox .text-center{text-align:center;}.qmbox .thumbnail{padding:4px;max-width:100%;border:1px solid #E3E9ED;}.qmbox .btn-primary{color:#fff;font-size:16px;padding:8px 14px;line-height:20px;border-radius:2px;display:inline-block;background:#FE7300;text-decoration:none;}.qmbox .btn-primary:hover,.qmbox .btn-primary:active{color:#fff;}.qmbox .footer{color:#9B9B9B;font-size:12px;margin-top:40px;}.qmbox .footer a{color:#9B9B9B;}.qmbox .footer a:hover{color:#fe9d4c;}.qmbox .footer a:active{color:#b15000;}.qmbox .email-body#mail_to_teacher{line-height:26px;color:#40485B;font-size:16px;padding:0px;}.qmbox .email-body#mail_to_teacher .container,.qmbox .email-body#mail_to_teacher .panel-body{padding:0px;}.qmbox .email-body#mail_to_teacher .container{padding-top:20px;}.qmbox .email-body#mail_to_teacher .textarea{padding:32px;}.qmbox .email-body#mail_to_teacher .say-hi{font-weight:500;}.qmbox .email-body#mail_to_teacher .paragraph{margin-top:24px;}.qmbox .email-body#mail_to_teacher .paragraph .pro-name{color:#000000;}.qmbox .email-body#mail_to_teacher .paragraph.link{margin-top:32px;text-align:center;}.qmbox .email-body#mail_to_teacher .paragraph.link .button{background:#4A90E2;border-radius:2px;color:#FFFFFF;text-decoration:none;padding:11px 17px;line-height:14px;display:inline-block;}.qmbox .email-body#mail_to_teacher ul.pro-desc{list-style-type:none;margin:0px;padding:0px;padding-left:16px;}.qmbox .email-body#mail_to_teacher ul.pro-desc li{position:relative;}.qmbox .email-body#mail_to_teacher ul.pro-desc li::before{content:'';width:3px;height:3px;border-radius:50%;background:red;position:absolute;left:-15px;top:11px;background:#40485B;}.qmbox .email-body#mail_to_teacher .blackboard-area{height:600px;padding:40px;background-image:url();color:#FFFFFF;}.qmbox .email-body#mail_to_teacher .blackboard-area .big-title{font-size:32px;line-height:45px;text-align:center;}.qmbox .email-body#mail_to_teacher .blackboard-area .desc{margin-top:8px;}.qmbox .email-body#mail_to_teacher .blackboard-area .desc p{margin:0px;text-align:center;line-height:28px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card:nth-child(odd){float:left;margin-top:45px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card:nth-child(even){float:right;margin-top:45px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card .title{font-size:18px;text-align:center;margin-bottom:10px;}\r\n" +
			"</style>\r\n" +
			"<meta>\r\n" +
			"<div class=\"email-body\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
			"<div class=\"container\">\r\n" +
			"<div class=\"logo\">\r\n" +
			"<img src=\"" + logo + "\",height=\"100\" width=\"100\">\r\n" +
			"</div>\r\n" +
			"<div class=\"panel\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
			"<div class=\"panel-header\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
			projectName + "邮件提醒\r\n" +
			"\r\n" +
			"</div>\r\n" +
			"<div class=\"panel-body\">\r\n" +
			"<p>" + text + "</p>\r\n" +
			"</div>\r\n" +
			"</div>\r\n" +
			"<div class=\"footer\">\r\n" +
			"<a href=\" " + dataWebsiteUrl + "\">@" + projectNameEn + "</a>\n" +
			"<div class=\"pull-right\"></div>\r\n" +
			"</div>\r\n" +
			"</div>\r\n" +
			"</div>\r\n" +
			"<style type=\"text/css\">.qmbox style, .qmbox script, .qmbox head, .qmbox link, .qmbox meta {display: none !important;}</style></div></div><!-- --><style>#mailContentContainer .txt {height:auto;}</style>  </div>\r\n" +
			" </body>\r\n" +
			"</html>"
	sendEmail(email, content)
}

func (*email) SendActiveEmail(user models.User, token string) {
	text :=
		"<html>\r\n" +
			" <head>\r\n" +
			"  <title>" + projectName + "</title>\r\n" +
			" </head>\r\n" +
			" <body>\r\n" +
			"  <div id=\"contentDiv\" onmouseover=\"getTop().stopPropagation(event);\" onclick=\"getTop().preSwapLink(event, 'spam', 'ZC1222-PrLAp4T0Z7Z7UUMYzqLkb8a');\" style=\"position:relative;font-size:14px;height:auto;padding:15px 15px 10px 15px;z-index:1;zoom:1;line-height:1.7;\" class=\"body\">    \r\n" +
			"  <div id=\"qm_con_body\"><div id=\"mailContentContainer\" class=\"qmbox qm_con_body_content qqmail_webmail_only\" style=\"\">\r\n" +
			"<style>\r\n" +
			"  .qmbox .email-body{color:#40485B;font-size:14px;font-family:-apple-system, \"Helvetica Neue\", Helvetica, \"Nimbus Sans L\", \"Segoe UI\", Arial, \"Liberation Sans\", \"PingFang SC\", \"Microsoft YaHei\", \"Hiragino Sans GB\", \"Wenquanyi Micro Hei\", \"WenQuanYi Zen Hei\", \"ST Heiti\", SimHei, \"WenQuanYi Zen Hei Sharp\", sans-serif;background:#f8f8f8;}.qmbox .pull-right{float:right;}.qmbox a{color:#FE7300;text-decoration:underline;}.qmbox a:hover{color:#fe9d4c;}.qmbox a:active{color:#b15000;}.qmbox .logo{text-align:center;margin-bottom:20px;}.qmbox .panel{background:#fff;border:1px solid #E3E9ED;margin-bottom:10px;}.qmbox .panel-header{font-size:18px;line-height:30px;padding:10px 20px;background:#fcfcfc;border-bottom:1px solid #E3E9ED;}.qmbox .panel-body{padding:20px;}.qmbox .container{width:50%;min-width:300px;padding:20px;margin:0 auto;}.qmbox .text-center{text-align:center;}.qmbox .thumbnail{padding:4px;max-width:100%;border:1px solid #E3E9ED;}.qmbox .btn-primary{color:#fff;font-size:16px;padding:8px 14px;line-height:20px;border-radius:2px;display:inline-block;background:#FE7300;text-decoration:none;}.qmbox .btn-primary:hover,.qmbox .btn-primary:active{color:#fff;}.qmbox .footer{color:#9B9B9B;font-size:12px;margin-top:40px;}.qmbox .footer a{color:#9B9B9B;}.qmbox .footer a:hover{color:#fe9d4c;}.qmbox .footer a:active{color:#b15000;}.qmbox .email-body#mail_to_teacher{line-height:26px;color:#40485B;font-size:16px;padding:0px;}.qmbox .email-body#mail_to_teacher .container,.qmbox .email-body#mail_to_teacher .panel-body{padding:0px;}.qmbox .email-body#mail_to_teacher .container{padding-top:20px;}.qmbox .email-body#mail_to_teacher .textarea{padding:32px;}.qmbox .email-body#mail_to_teacher .say-hi{font-weight:500;}.qmbox .email-body#mail_to_teacher .paragraph{margin-top:24px;}.qmbox .email-body#mail_to_teacher .paragraph .pro-name{color:#000000;}.qmbox .email-body#mail_to_teacher .paragraph.link{margin-top:32px;text-align:center;}.qmbox .email-body#mail_to_teacher .paragraph.link .button{background:#4A90E2;border-radius:2px;color:#FFFFFF;text-decoration:none;padding:11px 17px;line-height:14px;display:inline-block;}.qmbox .email-body#mail_to_teacher ul.pro-desc{list-style-type:none;margin:0px;padding:0px;padding-left:16px;}.qmbox .email-body#mail_to_teacher ul.pro-desc li{position:relative;}.qmbox .email-body#mail_to_teacher ul.pro-desc li::before{content:'';width:3px;height:3px;border-radius:50%;background:red;position:absolute;left:-15px;top:11px;background:#40485B;}.qmbox .email-body#mail_to_teacher .blackboard-area{height:600px;padding:40px;background-image:url();color:#FFFFFF;}.qmbox .email-body#mail_to_teacher .blackboard-area .big-title{font-size:32px;line-height:45px;text-align:center;}.qmbox .email-body#mail_to_teacher .blackboard-area .desc{margin-top:8px;}.qmbox .email-body#mail_to_teacher .blackboard-area .desc p{margin:0px;text-align:center;line-height:28px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card:nth-child(odd){float:left;margin-top:45px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card:nth-child(even){float:right;margin-top:45px;}.qmbox .email-body#mail_to_teacher .blackboard-area .card .title{font-size:18px;text-align:center;margin-bottom:10px;}\r\n" +
			"</style>\r\n" +
			"<meta>\r\n" +
			"<div class=\"email-body\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
			"<div class=\"container\">\r\n" +
			"<div class=\"logo\">\r\n" +
			"<img src=\"" + logo + "\",height=\"100\" width=\"100\">\r\n" +
			"</div>\r\n" +
			"<div class=\"panel\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
			"<div class=\"panel-header\" style=\"background-color: rgb(246, 244, 236);\">\r\n" +
			projectName + "账号激活\r\n" +
			"\r\n" +
			"</div>\r\n" +
			"<div class=\"panel-body\">\r\n" +
			"<p>您好 <a href=\"mailto:" + user.Email + "\" rel=\"noopener\" target=\"_blank\">" + user.NickName + "<wbr></a>！</p>\r\n" +
			"<p>欢迎您注册" + projectName + "账号，请点击下方链接进行账号激活</p>\r\n" +
			"<p>地址：" + "<a href=\"" + dataWebUrl + "/login/activeUser/" + token + "\">点击这里</a>" + "</p>\r\n" +
			"\r\n" +
			"</div>\r\n" +
			"</div>\r\n" +
			"<div class=\"footer\">\r\n" +
			"<a href=\" " + dataWebsiteUrl + "\">@" + projectNameEn + "</a>\n" +
			"<div class=\"pull-right\"></div>\r\n" +
			"</div>\r\n" +
			"</div>\r\n" +
			"</div>\r\n" +
			"<style type=\"text/css\">.qmbox style, .qmbox script, .qmbox head, .qmbox link, .qmbox meta {display: none !important;}</style></div></div><!-- --><style>#mailContentContainer .txt {height:auto;}</style>  </div>\r\n" +
			" </body>\r\n" +
			"</html>"
	sendEmail(user.Email, text)
}

var Email = &email{}
