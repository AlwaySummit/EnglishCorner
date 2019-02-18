package main

import (
	"time"
	"ec/english-corner/services"
	"ec/english-corner/common"
	"encoding/json"
	"fmt"
)

var currentTime = time.Unix(time.Now().Unix(), 0)
var attendTimeStr = currentTime.Format("2006-01-02")
var currentTimeUnix, _ = time.Parse("2006-01-02", attendTimeStr)

var contactInfo = map[string]string{}

func ResolveInfo(receivedFlag bool, attenderResp AttenderResp) ([]string) {
	attInfo := attenderResp.DateMember

	attenders := []string{}

	if err := json.Unmarshal([]byte(common.Credentials.PersonalInfo), &contactInfo); err != nil {
		fmt.Printf("[%s] resp to struct error: %+v", err)
	}

	for k, v := range attInfo {
		if currentTimeUnix.Unix() > attInfo[k].Timestamp {
			continue
		} else {
			// All hands
			if len(v.Members) == 1 {
				for k, _ := range contactInfo {
					attenders = append(attenders, k)
				}
			} else if receivedFlag {
				nonAttenderNames := []string{}
				// reverse the kv of contact map
				reverseContactInfo := reverseContactInfo(contactInfo)
				// get non attenders e.g. ["11111111"]
				nonAttenders := services.QueueWatcher()
				// map the corresponding names e.g. ["names"]
				for _, value := range nonAttenders {
					if v, ok := reverseContactInfo[value]; ok {
						nonAttenderNames = append(nonAttenderNames, v)
					}
				}

				// remove from the attenders
				for i := 0; i < len(v.Members); {
					flag := false
					for kk, _ := range nonAttenderNames {
						if v.Members[i] == nonAttenderNames[kk] {
							v.Members = remove(v.Members, i)
							flag = true
							break
						}
					}
					if !flag {
						i ++
					}
				}
				attenders = v.Members
			} else {
				attenders = v.Members
			}
			break
		}
		fmt.Printf("[%s] Attenders: %+v", attenders)
	}
	return attenders
}

func SendSms() (resp string, err error) {
	names := GetAttenderResults()
	attenders := ResolveInfo(false, names)
	numbers := services.ConstructNumbers(contactInfo, attenders)
	signName := services.ConstructSignName(attenders)
	templateParam := services.ConstructTemplateParam(false, attendTimeStr, attenders)
	resp, err = services.SendSMS(false, numbers, signName, templateParam)
	return resp, err
	//TODO err handlers
}

func SendReceiveSms() (resp string, err error) {
	names := GetAttenderResults()
	attenders := ResolveInfo(true, names)
	numbers := services.ConstructNumbers(contactInfo, attenders)
	signName := services.ConstructSignName(attenders)
	templateParam := services.ConstructTemplateParam(true, attendTimeStr, attenders)
	resp, err = services.SendSMS(true, numbers, signName, templateParam)
	return resp, err
	//TODO err handlers
}

func remove(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}

func reverseContactInfo(contactInfo map[string]string) (map[string]string) {
	reverseContactInfo := map[string]string{}
	for k, v := range contactInfo {
		reverseContactInfo[v] = k
	}
	return reverseContactInfo
}
