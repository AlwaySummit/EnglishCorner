package services

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"net/http"
	"encoding/json"
	"ec/english-corner/common"
)

type baseTemplateEle struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

func ConstructNumbers(dict map[string]string, names []string) (numbers string) {
	numbersArr := []string{}
	for _, value := range names {
		if v, ok := dict[value]; ok {
			numbersArr = append(numbersArr, v)
		}
	}
	numbersByte, _ := json.Marshal(numbersArr)
	numbers = string(numbersByte)
	return
}

func SendSMS(numbers string, signName string, templateParam string) {
	client, err := sdk.NewClientWithAccessKey(common.DEFAULT_REGION_ID, common.Credentials.AccessKey, common.Credentials.AccessKeySecret)
	if err != nil {
		fmt.Printf("[%s] Error when get client, err: %+v", err)
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Method = http.MethodPost
	request.Domain = common.SMS_DOMAIN
	request.Version = common.VERSION
	request.ApiName = common.API_NAME
	request.QueryParams["PhoneNumberJson"] = numbers
	request.QueryParams["SignNameJson"] = signName
	request.QueryParams["TemplateCode"] = common.Credentials.TemplateCode
	request.QueryParams["TemplateParamJson"] = templateParam

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		fmt.Printf("[%s] Error when get resp, err: %+v", err)
		panic(err)
	}
	fmt.Printf("[%s] Response of processing sms request", response.GetHttpContentString())
}

func ConstructSignName(names []string) (signName string) {
	baseSignName := common.Credentials.SignName
	signNameArr := []string{}
	for _, _ = range names {
		signNameArr = append(signNameArr, baseSignName)
	}

	signNameByte, _ := json.Marshal(signNameArr)
	signName = string(signNameByte)
	return
}

func ConstructTemplateParam(attendTimeStr string, names []string) (templateParam string) {
	var ele baseTemplateEle
	templateArr := []baseTemplateEle{}
	for _, v := range names {
		ele.Name = v
		ele.Time = attendTimeStr
		templateArr = append(templateArr, ele)
	}

	templateByte, _ := json.Marshal(templateArr)
	templateParam = string(templateByte)
	return
}
