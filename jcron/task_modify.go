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
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CronTask interface {
	New() bool
	Modify(id int, task *Task) bool
	Delete(id int) bool
	State(id int)
}

type Task struct {
	name string
	TaskFrequency
	command string
}

type TaskFrequency struct {
	Second string
	Minute string
	Hour string
	Day string
	Week string
	Month string
}

type TaskModify struct {
	operate string
	task *Task
}

type CronTime struct {
	cycleType int
	num int
}

const (
	fixedTime = iota
	cycleTime
	ignoreTime
)

var monthDayNumMap = [...]int{
	0 : 31,
	1 : 28,
	2 : 31,
	3 : 30,
	4 : 31,
	5 : 30,
	6 : 31,
	7 : 31,
	8 : 30,
	9 : 31,
	10 : 30,
	11: 31,
}

var (
	cronTable map[string]*Task
	tickChan map[string]chan time.Time
)

func (task *Task) New(fre *TaskFrequency) error {
	//用执行的命令作为任务唯一标识
	md5Str := task.command
	md5Obj := md5.New()
	md5Obj.Write([]byte(md5Str))
	md5Id := hex.EncodeToString(md5Obj.Sum(nil))
	cronTable[md5Id] = task
	tickSecond, err := getTickSecond(task)
	if err != nil {
		return err
	}
	go setTickCommand(tickSecond, task)
	return errors.New("Create error ")
}

func setTickCommand(second time.Duration, task *Task)  {
	tickChane := time.Tick(second)
	<-tickChane

}

func (task *Task) Modify(id int, newTask *Task) bool {
	return true
}

func (task *Task) Delete(id int) bool {
	return true
}

func (task *Task) State(id int) {

}

//秒 分 时 日 周 月
func getTickSecond(task *Task) (tickSecond time.Duration, err error) {
	var nextTime string
	nextTime += ""
	var second, minute, hour, day, week, month, year int
	week += 1
	year += 1
	now := time.Now()
	nowYear, nowMonth, nowDay := now.Date()

	var beyondDay int
	//月格式解析
	mon, err := parseTimeStr(task.month)
	if err != nil {
		return 0, err
	}
	switch mon.cycleType {
	case fixedTime:
		month = mon.num
	case ignoreTime:
		month = int(nowMonth)
	case cycleTime:
		month += int(nowMonth) + mon.num - 1
		if month > 12 {
			beyondDay += (month - 12) * 30
		}
	}

	//日格式解析
	d, err := parseTimeStr(task.week)
	if err != nil {
		return 0, err
	}
	//周格式解析
	we, err := parseTimeStr(task.week)
	if err != nil {
		return 0, err
	}
	if d.cycleType == ignoreTime {
		switch we.cycleType {
		case cycleTime:
			we.num += 1
			fallthrough
		case fixedTime:
			nowWeekday := now.Weekday()
			diff := we.num - int(nowWeekday)
			if diff < 0 {
				diff = diff + 7
			}
			day = nowDay + diff
		case ignoreTime:
		}
	}
	switch d.cycleType {
	case fixedTime:
		day = d.num
	case ignoreTime:
		day = nowDay
	case cycleTime:
		day = nowDay + d.num - 1
		//todo 超过当前月范围往月份上加
		thisMonthDayNum := monthDayNumMap[nowMonth]
		if nowMonth == 2 && nowYear / 4 == 0 {
			thisMonthDayNum ++
		}
		if day > thisMonthDayNum {
			day = 1
			beyondDay += day - thisMonthDayNum
		}
	}


	//秒格式解析
	s, err := parseTimeStr(task.second)
	if err != nil {
		return 0, err
	}
	if s.cycleType == fixedTime {
		second += s.num
	} else {
		second = now.Second() + s.num
		if second >= 60 {
			minute += second / 60
			second = second % 60
		}
	}

	//分钟格式解析
	min, err := parseTimeStr(task.minute)
	if err != nil {
		return 0, err
	}
	if min.cycleType == fixedTime {
		minute += min.num
	} else {
		minute += now.Minute() + min.num
		if minute >= 60 {
			hour += minute / 60
			minute = minute % 60
		}
	}
	h, err := parseTimeStr(task.hour)
	if err != nil {
		return 0, err
	}
	if h.cycleType == fixedTime {
		hour += h.num
	} else {
		hour = now.Hour() + h.num
		if hour >= 24 {
			day += hour / 24
			minute = hour % 24
		}
	}
	d, err = parseTimeStr(task.day)
	if err != nil {
		return 0, nil
	}
	if d.cycleType == fixedTime {
		day += d.num
	} else {
		nowYear, nowMonth, nowDay := now.Date()
		thisMonthDayNum := monthDayNumMap[nowMonth]
		if nowMonth == 2 && nowYear / 4 == 0 {
			thisMonthDayNum ++
		}
		day += nowDay
		if day >= thisMonthDayNum{
		}
	}

	return time.Duration(second), nil
}

func parseTimeStr(str string) (CronTime, error) {
	if num, err := strconv.Atoi(str); err == nil && num != 0 {
		//为数字则为固定时间
		return CronTime{
			cycleType: fixedTime,
			num: num,
		}, nil
	}
	if str == "*" {
		return CronTime{
			cycleType: ignoreTime,
			num: 1,
		}, nil
	} else {
		regexpObj, _ := regexp.Compile(`[1-9][0-9]+/\*?`)
		timStr := regexpObj.FindString(str)
		if timStr == "" {
			 log.Println("Time format error : " + str)
			 return CronTime{}, errors.New("Time format error : " + str)
		}
		arr := strings.Split("/", timStr)
		timeInt, _ := strconv.Atoi(arr[0])
		return CronTime{
			cycleType:cycleTime,
			num:timeInt,
		}, nil
	}
}

