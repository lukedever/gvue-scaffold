package helper

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dm"
	"github.com/lukedever/gvue-scaffold/internal/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func SendEmail(to, sub, body string) error {
	accessKeyId := viper.GetString("email.id")
	accessKeySecret := viper.GetString("email.secret")
	if accessKeyId == "" || accessKeySecret == "" {
		//TODO:记录日志
		log.Warn("email key is empty")
		return nil
	}
	client, err := dm.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	request := dm.CreateSingleSendMailRequest()
	request.Scheme = "https"
	request.AccountName = viper.GetString("email.from")
	request.AddressType = requests.NewInteger(1)
	request.ReplyToAddress = requests.NewBoolean(false)
	request.ToAddress = to
	request.Subject = sub
	request.FromAlias = viper.GetString("app.name")
	request.HtmlBody = body
	request.ReplyToAddress = requests.NewBoolean(true)

	response, err := client.SingleSendMail(request)
	if err != nil {
		log.Error("send email error", zap.Error(err))
		return err
	}
	log.Debug("response is: " + response.String())
	return nil
}