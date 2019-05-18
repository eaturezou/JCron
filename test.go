/*
 | ---------------------------------------------------------
 | Author: Zoueature
 | Email: zoueature@gmail.com
 | Date: 2019/5/13
 | Time: 9:05
 | Description:
 | ---------------------------------------------------------
*/

package main

import (
	"fmt"
	"learning/JCron/jcron"
)

func main() {
	task := jcron.Task{
		Name: "hello world",
		TaskFrequency: jcron.TaskFrequency{
			Second:"1",
			Minute:"2",
			Hour:"*",
			Day:"2/*",
			Month:"*",
			Week:"*",
		},
		Command:"php -r 'echo 123;'",
	}
	timestamp, _ := jcron.GetTickSecond(&task)
	cronTask := jcron.CronTask{
		Id:"1892213109",
		ExecuteTime:timestamp,
	}
	cronTask2 := jcron.CronTask{
		Id:"chdsodfoia",
		ExecuteTime:321313131,
	}
	queue := &jcron.CronTask{}
	queue.Insert(&cronTask)
	queue.Insert(&cronTask2)
	for {
		fmt.Printf("%+v\n", queue)
		if queue.Next != nil {
			queue = queue.Next
		} else {
			break
		}
	}
}