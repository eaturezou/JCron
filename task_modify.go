/* -------------------------------------------------
| Author: Zoueature
| Email: zoueature@gmail.com
| Date: 19-5-12
| Description: 
| -------------------------------------------------
*/

package main

import "net"

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

func (task *Task) New(fre *TaskFrequency) {

}

func doTaskModify(conn net.Conn) {

}


func parseQuery(conn net.Conn) {

}