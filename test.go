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
	"time"
)

func main() {
	testTask()
	//testLink()

}

var cronQueue = &jcron.CronTask{}

func testLink() {
	cronTask := &jcron.CronTask{
		Id:"1",
		ExecuteTime:1,
	}
	cronTask2 := &jcron.CronTask{
		Id:"2",
		ExecuteTime:2,
	}
	cronTask3 := &jcron.CronTask{
		Id:"3",
		ExecuteTime:3,
	}
	cronQueue.Insert(cronTask)
	cronQueue.Insert(cronTask2)
	cronQueue.Insert(cronTask3)
	//printLink(cronQueue)
	_, _ = cronQueue.Delete("1")
	_, _ = cronQueue.Delete("2")
	_, _ = cronQueue.Delete("3")
	printLink(cronQueue)
}

func testTask()  {
	task := &jcron.Task{
		Name: "hello world",
		TaskFrequency: jcron.TaskFrequency{
			Second:"10/*",
			Minute:"*",
			Hour:"*",
			Day:"*",
			Month:"*",
			Week:"*",
		},
		Command:"php -r 'echo 123;'",
	}
	err := jcron.New(task)
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		fmt.Println("------------------Heart Beat------------------")
		time.Sleep(5*time.Second)
	}
}

func printLink(node *jcron.CronTask) {
	for {
		if node != nil {
			fmt.Printf("%+v\n", node)
		}
		if node.Next == nil {
			break
		} else {
			node = node.Next
		}
	}
}