/* -------------------------------------------------
| Author: Zoueature
| Email: zoueature@gmail.com
| Date: 19-5-12
| Description:
| -------------------------------------------------
*/

package jcron

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	Name string
	TaskFrequency
	Command string
}

type TaskFrequency struct {
	Second string
	Minute string
	Hour   string
	Day    string
	Week   string
	Month  string
}

type CronTime struct {
	cycleType int
	num       int
}

type TaskCommand struct {
	TaskType int
	content  string
}

const (
	script = iota
	httpCallback
)

const (
	fixedTime = iota
	cycleTime
	ignoreTime
)

var monthDayNumMap = [...]int{
	0:  31,
	1:  28,
	2:  31,
	3:  30,
	4:  31,
	5:  30,
	6:  31,
	7:  31,
	8:  30,
	9:  31,
	10: 30,
	11: 31,
}

var (
	cronQueue = &CronTask{}
	mutex     sync.RWMutex
)

func init() {
	go dispatcher()
}

//秒 分 时 日 周 月
//获取下次执行的时间， 字符串日期格式
func GetTickSecond(task *Task) (tickSecond int64, err error) {
	nowTime := time.Now()
	factYear, nowMonth, factDay := nowTime.Date()
	factMonth := int(nowMonth)
	factHour := nowTime.Hour()
	factMinute := nowTime.Minute()
	factSecond := nowTime.Second()
	secondTimeCron, err := parseTimeStr(task.Second)
	if err != nil {
		return 0, err
	}
	minuteTimeCron, err := parseTimeStr(task.Minute)
	if err != nil {
		return 0, err
	}
	hourCronTime, err := parseTimeStr(task.Hour)
	if err != nil {
		return 0, err
	}
	dayCronTime, err := parseTimeStr(task.Day)
	if err != nil {
		return 0, err
	}
	monthCronTime, err := parseTimeStr(task.Month)
	if err != nil {
		return 0, err
	}
	factSecond = parseCronTime(&secondTimeCron, factSecond)
	if secondTimeCron.cycleType != fixedTime {
		goto spliceTime
	}
	factMinute = parseCronTime(&minuteTimeCron, factMinute)
	if minuteTimeCron.cycleType != fixedTime {
		goto spliceTime
	}
	factHour = parseCronTime(&hourCronTime, factHour)
	if minuteTimeCron.cycleType != fixedTime {
		goto spliceTime
	}
	factDay = parseCronTime(&dayCronTime, factDay)
	if dayCronTime.cycleType != fixedTime {
		goto spliceTime
	}
	factMonth = parseCronTime(&monthCronTime, factMonth)
spliceTime:
	if factSecond >= 60 {
		factMinute += factSecond / 60
		factSecond = factSecond % 60
	}
	if factMinute >= 60 {
		factHour += factMinute / 60
		factMinute = factMinute % 60
	}
	if factHour >= 24 {
		factDay += factHour / 24
		factHour = factHour % 24
	}
	thisDayNum := monthDayNumMap[nowMonth]
	if factDay > thisDayNum {
		for i := nowMonth; i < 12; i++ {
			if factDay > monthDayNumMap[i] {
				factDay -= monthDayNumMap[i]
				factMonth++
			} else {
				break
			}
		}
	}
	if factMonth > 12 {
		factYear = factYear + factMonth/12
		factMonth = factMonth % 12
	}
	factDate := fmt.Sprintf("%02d", factYear) + "-" +
		fmt.Sprintf("%02d", factMonth) + "-" +
		fmt.Sprintf("%02d", factDay) + " " +
		fmt.Sprintf("%02d", factHour) + ":" +
		fmt.Sprintf("%02d", factMinute) + ":" +
		fmt.Sprintf("%02d", factSecond)
	loc, _ := time.LoadLocation("Local")
	timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", factDate, loc)
	if err != nil {
		return -1, err
	}
	result := timestamp.Unix()
	return result, nil
}

/*
解析实际的执行时间
*/
func parseCronTime(cronTime *CronTime, referTime int) (fact int) {
	switch cronTime.cycleType {
	case ignoreTime:
		return referTime
	case fixedTime:
		return cronTime.num
	case cycleTime:
		return cronTime.num + referTime
	}
	return referTime
}

func parseTimeStr(str string) (CronTime, error) {
	if num, err := strconv.Atoi(str); err == nil && num != 0 {
		//为数字则为固定时间
		return CronTime{
			cycleType: fixedTime,
			num:       num,
		}, nil
	}
	if str == "*" {
		return CronTime{
			cycleType: ignoreTime,
			num:       1,
		}, nil
	} else {
		regexpObj, _ := regexp.Compile(`[1-9][0-9]*/\*?`)
		timStr := regexpObj.FindString(str)
		if timStr == "" {
			log.Println("Time format error : " + str)
			return CronTime{}, errors.New("Time format error : " + str)
		}
		arr := strings.Split(timStr, "/")
		timeInt, _ := strconv.Atoi(arr[0])
		return CronTime{
			cycleType: cycleTime,
			num:       timeInt,
		}, nil
	}
}

//新增一个任务的任务队列
func New(task *Task) error {
	//用执行的命令作为任务唯一标识
	md5Str := task.Command
	md5Obj := md5.New()
	md5Obj.Write([]byte(md5Str))
	md5Id := hex.EncodeToString(md5Obj.Sum(nil))
	timestamp, err := GetTickSecond(task)
	if timestamp <= 0 {
		return errors.New("Parse times error, " + err.Error())
	}
	if err != nil {
		return err
	}
	nowTimestamp := time.Now().Unix()
	if nowTimestamp == timestamp {
		timestamp++
	}
	newCronTsk := &CronTask{
		Id:          md5Id,
		ExecuteTime: timestamp,
		Task:        task,
	}
	cronQueue.Insert(newCronTsk)
	return nil
}

func GetTask() *CronTask {
	mutex.Lock()
	defer mutex.Unlock()
	return cronQueue.GetFirst()
}

func DeleteTask(id string, lock bool) error {
	if lock {
		mutex.Lock()
		defer mutex.Unlock()
	}
	_, err := cronQueue.Delete(id)
	return err
}
