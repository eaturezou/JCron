/*
 | ---------------------------------------------------------
 | Author: Zoueature
 | Email: zoueature@gmail.com
 | Date: 2019/5/16
 | Time: 11:05
 | Description:
 | ---------------------------------------------------------
*/

package jcron

import (
	"errors"
)

type CronQueue interface {
	Get(id string) (interface{}, error)
	Delete(id string) (interface{}, error)
	Insert(task *CronTask) interface{}
	GetFirst() interface{}
}

type CronTask struct {
	Id string
	ExecuteTime int64
	Task *Task
	Next *CronTask
	Prev *CronTask
}
//1,3
func (node *CronTask) Insert(task *CronTask) *CronTask {
	executeTime := task.ExecuteTime
	for {
		if node.ExecuteTime >= executeTime {
			task.Prev = node.Prev
			task.Next = node
			node.Prev.Next = task
			node.Prev = task
			break
		} else {
			if node.Next == nil {
				node.Next = task
				task.Prev = node
				break
			}
			node = node.Next
		}
	}
	return node
}

func (node *CronTask) Get(id string) (*CronTask, error) {
	for {
		if node.Id == id {
			return node, nil
		} else {
			if node.Next != nil {
				break
			} else {
				node = node.Next
			}
		}
	}
	return nil, errors.New("No Result ")
}

func (node *CronTask) GetFirst() *CronTask {
	if node.Next == nil {
		return nil
	}
	return node.Next
}

func (node *CronTask) Delete(id string) (*CronTask, error) {
	if node == nil || node.Next == nil {
		return node, errors.New("empty queue ")
	}
	for {
		if node.Id == id {
			if node.Prev != nil {
				node.Prev.Next = node.Next
			}
			if node.Next != nil {
				node.Next.Prev = node.Prev
			}
			return node, nil
		} else {
			if node.Next != nil {
				node = node.Next
			} else {
				break
			}
		}
	}
	return nil, errors.New("Not Found ")
}