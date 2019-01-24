package common

import (
	"os"
)

type Configuration struct {
	AccessKey       string
	AccessKeySecret string
	PersonalInfo    string
	Region          string
	SignName        string
	SmsQueueName    string
	TemplateCode    string
}

var Credentials Configuration

func Retrieve() {
	Credentials.AccessKey = os.Getenv(ACCESS_KEY)
	Credentials.AccessKeySecret = os.Getenv(ACCESS_SECRET)
	Credentials.SignName = os.Getenv(SIGN_NAME)
	Credentials.Region = os.Getenv(REGION_ID)
	Credentials.SmsQueueName = os.Getenv(SMS_QUEUE_NAME)
	Credentials.TemplateCode = os.Getenv(TEMPLATE_CODE)
	Credentials.PersonalInfo = os.Getenv(PERSONAL_INFO)
}
