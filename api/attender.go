package main

import (
	"time"
)

type AttenderInfo struct {
	Members   []string `json:"members"`
	Timestamp int64    `json:"timestamp"`
	Date      string   `json:"date"`
	WeekDay   string   `json:"day"`
	Words     []string `json:"words"`
}

type AttenderResp struct {
	DateMember []AttenderInfo `json:"date_member"`
	LatestDate int64          `json:"latest_date"`
}

func GetAttenders() (result [][]string) {
	array := []string{"Marco", "Craig", "Patrick", "Mike", "Jonas", "Fred", "Christina"}
	for n := 0; n < 300; n++ {
		group1 := []string{}
		group2 := []string{}
		group3 := []string{}
		group1 = append(group1, array[0], array[1], array[2], array[3], array[4])
		group2 = append(group2, array[2], array[3], array[4], array[5], array[6])
		group3 = append(group3, array[0], array[1], array[5], array[6], "Min")
		result = append(result, group1, group2, group3)
		array = []string{array[1], array[2], array[3], array[4], array[5], array[6], array[0]}
	}
	return
}

func GetAllDates() (dates []time.Time) {
	baseDate := time.Date(2019, 1, 11, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 300; i++ {
		dayOfWeek := baseDate.Weekday()
		if (dayOfWeek == 1 || dayOfWeek == 3 || dayOfWeek == 4) {
			dates = append(dates, baseDate)
		}
		baseDate = baseDate.AddDate(0, 0, 1)
	}
	return
}

func GetAttenderResults() (attenderResp AttenderResp) {
	attenders := GetAttenders()
	currentTime := time.Now()
	allDates := GetAllDates()
	dayOfWeek := []string{}
	participators := [][]string{}
	occurTimeObj := []time.Time{}
	var latestDate int64
	flag := true
	var j = 0
	for i := 0; i+2 < len(allDates); i++ {
		currentOccurTime := allDates[i]
		nextOccurTime := allDates[i+1]
		nextOfNextOccurTime := allDates[i+2]
		currentOccurDayOfWeek := currentOccurTime.Weekday().String()
		nextOccurDayOfWeek := nextOccurTime.Weekday().String()
		allHands := []string{"All Hands"}
		if (currentOccurTime.Month() != nextOfNextOccurTime.Month()) {
			participators = append(participators, allHands, allHands)
			dayOfWeek = append(dayOfWeek, currentOccurDayOfWeek, nextOccurDayOfWeek)
			occurTimeObj = append(occurTimeObj, currentOccurTime, nextOccurTime)
			i = i + 1
			continue
		}

		if currentOccurTime.YearDay() >= currentTime.YearDay() && flag {
			latestDate = currentOccurTime.Unix()
			flag = false
		}
		participators = append(participators, attenders[j])
		dayOfWeek = append(dayOfWeek, currentOccurDayOfWeek)
		occurTimeObj = append(occurTimeObj, currentOccurTime)
		j++

	}

	results := AttenderInfo{}
	dateMembers := []AttenderInfo{}
	for k, _ := range occurTimeObj {
		results.Timestamp = occurTimeObj[k].Unix()
		results.Date = time.Unix(results.Timestamp, 0).Format("2006-01-02")
		results.Members = participators[k]
		results.WeekDay = occurTimeObj[k].Weekday().String()
		dateMembers = append(dateMembers, results)
	}
	attenderResp.DateMember = dateMembers
	attenderResp.LatestDate = latestDate
	return
}
