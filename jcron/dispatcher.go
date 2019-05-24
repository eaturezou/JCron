/*
 | ---------------------------------------------------------
 | Author: Zoueature
 | Email: zoueature@gmail.com
 | Date: 2019/5/13
 | Time: 20:23
 | Description:
 | ---------------------------------------------------------
*/

package jcron

import (
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"
)

type Queue CronTask

var runResult chan bool

func dispatcher() {
	for {
		nowTimestamp := time.Now().Unix()
		task := GetTask()
		if task == nil {
			continue
		}
		diffSeconds := task.ExecuteTime - nowTimestamp
		if diffSeconds <= 0 {
			executeCommand(task)
			continue
		}
		tickChan := time.Tick(time.Second * time.Duration(diffSeconds))
		select {
		case <-tickChan:
			executeCommand(task)
		case result := <-runResult:
			if result {
				log.Println("Run success ")
			} else {
				log.Println("Run command fail ")
			}
		}
	}
}

func executeCommand(task *CronTask)  {
	mutex.Lock()
	defer mutex.Unlock()
	if task == nil || task.Task == nil {
		runResult<-false
		return
	}
	reg, _ := regexp.Compile(`^http(s)?://.*`)
	matched := reg.Match([]byte(task.Task.Command))
	go func() {
		if matched {
			//http回调
			response, err := http.Get(task.Task.Command)
			if err != nil && response.StatusCode == http.StatusOK {
				runResult<-true
			} else {
				runResult<-false
			}
		} else {
			//系统下脚本
			cmd := exec.Command("/usr/local/sbin/php", "-r", "'echo 123;'")
			msg,err := cmd.Output()
			log.Println(msg)
			if err != nil {
				runResult<-false
			} else {
				runResult<-true
			}
		}
	}()
	err := DeleteTask(task.Id)
	if err != nil {
		log.Println(err.Error())
	}

	err = New(task.Task)
	if err != nil {
		log.Println(err.Error())
	}
}
