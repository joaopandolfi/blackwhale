package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var coreDate time.Time

func DateSearchFormat(day string, month string, year string) string {
	return year + month + day
}

func DateToString(date time.Time) (year string, month string, day string) {
	strDate := date.Format("2006-01-02")
	splited := strings.Split(strDate, "-")
	year = splited[0]
	month = splited[1]
	day = splited[2]
	return
}

func DateToInt(date time.Time) int{
	strDate := date.Format("20060102")
	idate,_ := strconv.Atoi(strDate)
	return idate
}

func TimestampToString(date int) (year string, month string, day string){
	y := date/10000
	date = date - (y*10000)
	m := date/100
	d := date - (m*100)

	year = strconv.Itoa(y)
	month = strconv.Itoa(m)
	day = strconv.Itoa(d)

	if(m < 10){
		month = fmt.Sprintf("0%s",month)
	}
	if(d < 10){
		day = fmt.Sprintf("0%s",day)
	}

	return
}

func TimestampToIso(date int) (isoDate string){
	y,m,d := TimestampToString(date)
	return fmt.Sprintf("%s-%s-%s",y,m,d)
}

func TimestampToDate(timestamp int) string{
	y,m,d := TimestampToString(timestamp)
	return fmt.Sprintf("%s/%s/%s",d,m,y)
}

// Formar dd/mm/yyyy
func StringToTimestamp(stringTimestamp string) int{
	splited := strings.Split(stringTimestamp, "/")
	year,_ := strconv.Atoi(splited[2])
	month,_ := strconv.Atoi(splited[1])
	day,_ := strconv.Atoi(splited[0])

	return (year*10000)+(month*100)+(day)
}

func IsoToTimestamp(stringTimestamp string) int{
	splited := strings.Split(stringTimestamp, "-")
	year,_ := strconv.Atoi(splited[0])
	month,_ := strconv.Atoi(splited[1])
	day,_ := strconv.Atoi(splited[2])

	return (year*10000)+(month*100)+(day)
}

func IsoToTime(stringTimestamp string) (t time.Time){
	//splited := strings.Split(stringTimestamp, "T")

	t,_ = time.Parse("2006-01-02T15:04:05",stringTimestamp)
	return
}

func CurrentDate() time.Time {
	return time.Now().UTC()
}

func SetCoreDateTimestamp(d int){
	strDate := TimestampToIso(d)
	coreDate,_ = time.Parse("2006-01-02",strDate)
}

func SetCoreDate(d time.Time){
	coreDate = d
}

func GetCoreDate() time.Time{
	return coreDate
}

func GetCurrentISOTimestamp() string{
	return time.Now().UTC().Format("2006-01-02T15:04:05")
}

func GetlastDayMonth(month string) int {
	mt,_ := strconv.Atoi(month)
	mt--
	days := [13]int{31,28,31,30,31,30,31,31,30,31,30,31,-1}
	if(mt >=12 ){
		mt = 12
	}
	return days[mt]
}

func LastDayTimestamp() int {
	t := time.Now()
	t = t.AddDate(0,0,-1)
	timestamp, _ := strconv.Atoi(t.Format("20060102")) // 20060102150405 -> AAMMDDhhmmss
	return timestamp
}

func WindowedDate(pastWindow int, centerDate time.Time, fowardWindow int) (dates []time.Time) {

	deltaToZeroDay := (centerDate.Day() - 1) * -1
	for i := pastWindow; i > 0; i-- {
		dates = append(dates, centerDate.AddDate(0, (i*-1), deltaToZeroDay))
	}

	dates = append(dates, centerDate)

	for i := 1; i <= fowardWindow; i++ {
		dates = append(dates, centerDate.AddDate(0, i, deltaToZeroDay))
	}
	return
}
