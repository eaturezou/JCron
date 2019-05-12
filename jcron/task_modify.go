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
	second string
	minute string
	hour string
	day string
	week string
	month string
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
)

var cronTable map[string]*Task

func (task *Task) New(fre *TaskFrequency) bool {
	md5Str := task.name + task.second + task.minute + task.hour + task.day + task.week + task.month + task.command
	md5Obj := md5.New()
	md5Obj.Write([]byte(md5Str))
	md5Id := hex.EncodeToString(md5Obj.Sum(nil))
	cronTable[md5Id] = task
	return true
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
func getTickSecond(task *Task) (time.Duration, error) {
	var second time.Duration
	sec, err := parseTimeStr(task.second)
	if err != nil {
		return 0, nil
	}
	second += sec
	min, err := parseTimeStr(task.minute)
	if err != nil {
		return 0, nil
	}
	second += 60 * min

	return second, nil
}

func parseTimeStr(str string) (duration time.Duration, err error) {
	if str == "*" {
		return 1, nil
	} else {
		regexp, _ := regexp.Compile("[1-9][0-9]+/\*?")
		timStr := regexp.FindString(str)
		if timStr == "" {
			 log.Println("Time format error : " + str)
			 return 0, errors.New("Time format error : " + str)
		}
		arr := strings.Split("/", timStr)
		timeInt, _ := strconv.Atoi(arr[0])
		return time.Duration(timeInt), nil
	}
}