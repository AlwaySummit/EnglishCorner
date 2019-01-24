package main

import (
	"time"
	"ec/english-corner/services"
	"ec/english-corner/common"
	"encoding/json"
	"fmt"
)

var currentTimeUnix = time.Now().Unix() * 1000
var currentTime = time.Unix(currentTimeUnix/1000, 0)
var attendTimeStr = currentTime.Format("2006-01-02")

var contactInfo = map[string]string{}

func ResolveInfo(attenderResp AttenderResp) ([]string) {
	attInfo := attenderResp.DateMember

	attenders := []string{}

	if err := json.Unmarshal([]byte(common.Credentials.PersonalInfo), &contactInfo); err != nil {
		fmt.Printf("[%s] resp to struct error: %+v", err)
	}

	for k, v := range attInfo {
		if currentTimeUnix > attInfo[k].Timestamp {
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
		fmt.Printf("[%s] Attenders: %+v", attenders)
	}
	return attenders
}

func SendSms() (attenderResp AttenderResp) {
	names := GetAttenderResults()
	attenders := ResolveInfo(names)
	numbers := services.ConstructNumbers(contactInfo, attenders)
	signName := services.ConstructSignName(attenders)
	templateParam := services.ConstructTemplateParam(attendTimeStr, attenders)
	services.SendSMS(numbers, signName, templateParam)
	return attenderResp
	//TODO err handlers
}