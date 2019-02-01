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
}

var Credentials *Configuration

func Retrieve() {
	Credentials = &Configuration{
		AccessKey:       os.Getenv(ACCESS_KEY),
		AccessKeySecret: os.Getenv(ACCESS_SECRET),
		SignName:        os.Getenv(SIGN_NAME),
		Region:          os.Getenv(REGION_ID),
		SmsQueueName:    os.Getenv(SMS_QUEUE_NAME),
		PersonalInfo:    os.Getenv(PERSONAL_INFO),
	}
}
