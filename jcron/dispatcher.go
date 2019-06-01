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
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"
)

type Queue CronTask

var (
	runResult = make(chan bool)
	tickTask = make(chan *CronTask)
)

func dispatcher() {
	go func() {
		for {
			select {
			case task := <- tickTask:
				executeCommand(task)
			case result := <-runResult:
				if result {
					log.Println("Run success ")
				} else {
					log.Println("Run command fail ")
				}
			default:

			}
		}
	}()
	for {
		nowTimestamp := time.Now().Unix()
		task := GetTask()
		if task == nil {
			continue
		}
		diffSeconds := task.ExecuteTime - nowTimestamp
		if diffSeconds <= 0 {
			tickTask<-task
			continue
		}
		ticker := time.NewTicker(time.Second * time.Duration(diffSeconds))
		timestampMutex.Lock()
		nowExecuteTimestamp = task.ExecuteTime
		timestampMutex.Unlock()
		select {
		case <-ticker.C:
			tickTask<-task
		case <-hadInsert:
			ticker.Stop()
		}
	}
}

func executeCommand(task *CronTask) {
	mutex.Lock()
	defer mutex.Unlock()
	if task == nil || task.Task == nil {
		runResult <- false
		return
	}
	reg, _ := regexp.Compile(`^http(s)?://.*`)
	matched := reg.Match([]byte(task.Task.Command))
	go func() {
		if matched {
			//http回调
			response, err := http.Get(task.Task.Command)
			_, err = ioutil.ReadAll(response.Body)
			if response != nil && response.Body != nil {
				defer response.Body.Close()
			}
			if err != nil && response.StatusCode == http.StatusOK {
				runResult <- true
			} else {
				runResult <- false
			}
		} else {
			//系统下脚本
			cmd := exec.Command("/usr/local/sbin/php", "-r", "'echo 123;'")
			_, err := cmd.Output()
			if err != nil {
				runResult <- false
			} else {
				runResult <- true
			}
		}
	}()
	err := DeleteTask(task.Id, false)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = New(task.Task, true)
	if err != nil && err.Error() != "expire " {
		log.Println(err.Error())
	}
}
