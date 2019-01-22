package services

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"ec/english-corner/common"
)

var currentTimeUnix = time.Now().Unix() * 1000
var currentTime = time.Unix(currentTimeUnix/1000, 0)
var logTimeStr = currentTime.Format("2006-01-02 03:04:05 PM")
var attendTimeStr = currentTime.Format("2006-01-02")

var contactInfo = map[string]string{
	"Craig":     "18966808305",
	"Patrick":   "15829603595",
	"Fred":      "15929934999",
	"Marco":     "18891993567",
	"Christina": "13152006719",
	"Min":       "18629092397",
	"Jonas":     "18710849386",
	"Mike":      "18089201617",
}

type baseTemplateEle struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

type respBody struct {
	Message    string  `json:"msg"`
	Result     respStr `json:"result"`
	ReturnCode string  `json:"retCode"`
}

type respStr struct {
	Value string `json:"v"`
}

type attendanceInfo struct {
	Date    int64    `json:"date"`
	Members []string `json:"members"`
}

func constructNumbers(dict map[string]string, names []string) (numbers string) {
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

func getAttenders(url string) []string {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[%s] Error when request for attenders, err: %+v", logTimeStr, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[%s] Error when read body, err: %+v", logTimeStr, err)
	}
	jsonStr := string(body)
	fmt.Printf("[%s] Attenders info from kv: %+v", logTimeStr, jsonStr)
	var dat respBody
	if err := json.Unmarshal([]byte(jsonStr), &dat); err != nil {
		fmt.Printf("[%s] resp to struct error: %+v", logTimeStr, err)
	}

	attInfo := []attendanceInfo{}
	if err := json.Unmarshal([]byte(dat.Result.Value), &attInfo); err != nil {
		fmt.Printf("[%s] Attenders info to struct error: %+v", logTimeStr, err)
	}

	attenders := []string{}

	for k, v := range attInfo {
		if currentTimeUnix > attInfo[k].Date {
			continue
		} else {
			// All hands
			if len(v.Members) == 1 {
				for k, _ := range contactInfo {
					attenders = append(attenders, k)
				}
			} else {
				attenders = v.Members
			}
			break
		}
		fmt.Printf("[%s] Attenders: %+v", logTimeStr, attenders)
	}
	return attenders
}

func sendSMS(numbers string, signName string, templateParam string) {
	client, err := sdk.NewClientWithAccessKey(common.DEFAULT_REGION_ID, common.ACCESS_KEY, common.ACCESS_SECRET)
	if err != nil {
		fmt.Printf("[%s] Error when get client, err: %+v", logTimeStr, err)
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Method = http.MethodPost
	request.Domain = common.SMS_DOMAIN
	request.Version = common.VERSION
	request.ApiName = common.API_NAME
	request.QueryParams["PhoneNumberJson"] = numbers
	request.QueryParams["SignNameJson"] = signName
	request.QueryParams["TemplateCode"] = common.TEMPLATE_CODE
	request.QueryParams["TemplateParamJson"] = templateParam

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		fmt.Printf("[%s] Error when get resp, err: %+v", logTimeStr, err)
		panic(err)
	}
	fmt.Printf("[%s] Response of processing sms request", logTimeStr, response.GetHttpContentString())
}

func constructSignName(names []string) (signName string) {
	baseSignName := common.SIGN_NAME
	signNameArr := []string{}
	for _, _ = range names {
		signNameArr = append(signNameArr, baseSignName)
	}

	signNameByte, _ := json.Marshal(signNameArr)
	signName = string(signNameByte)
	return
}

func constructTemplateParam(names []string) (templateParam string) {
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

func main() {
	names := getAttenders(URL)
	numbers := constructNumbers(contactInfo, names)
	signName := constructSignName(names)
	templateParam := constructTemplateParam(names)
	sendSMS(numbers, signName, templateParam)
}
