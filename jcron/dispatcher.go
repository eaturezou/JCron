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
	"fmt"
	"log"
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
		fmt.Println("Get task " + task.Id)
		diffSeconds := task.ExecuteTime - nowTimestamp
		if diffSeconds <= 0 {
			go executeCommand(task)
			continue
		}
		tickChan := time.Tick(time.Second * time.Duration(diffSeconds))
		select {
		case <-tickChan:
			go executeCommand(task)
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
	go func() {
		//cmd := exec.Command(task.Task.Command)
		//err:= cmd.Run()
		fmt.Println("Do the command ")
		var err error
		if err != nil {
			runResult<-false
		} else {
			runResult<-true
		}
	}()
	_, err := cronQueue.Delete(task.Id)
	if err != nil {
		log.Println(err.Error())
	}

}
