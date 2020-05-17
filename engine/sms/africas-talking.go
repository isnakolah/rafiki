package sms

import (
	"fmt"
	"github.com/AndroidStudyOpenSource/africastalking-go/sms"
	. "rafiki/settings"
)

var env string

func SendATMessage(recipientNumber string, messageBody string) (string, bool) {

	if GetEnv() == "PRODUCTION" {
		env = "Production"
	} else {
		env = "Production"
	}

	smsService := sms.NewService(GetAfricasTalkingUsername(), GetAfricasTalkingKey(), env)

	smsResponse, err := smsService.Send("", recipientNumber, messageBody)
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	fmt.Println(smsResponse)
	return smsResponse.SMS.Recipients[0].MessageID, true
}
