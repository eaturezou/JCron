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
	Next *CronTask
	Prev *CronTask
}
//1,3
func (node *CronTask) Insert(task *CronTask) *CronTask {
	if node.Id == "" {
		node.Id = task.Id
		node.ExecuteTime = task.ExecuteTime
		return node
	}
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
	for node.Next != nil  {
		if node.Id == id {
			return node, nil
		} else {
			node = node.Next
		}
	}
	return nil, errors.New("No Result ")
}

func (node *CronTask) GetFirst() *CronTask {
	return node
}

func (node *CronTask) Delete(id string) (*CronTask, error) {
	res := node
	for node.Next != nil  {
		if node.Id == id {
			node.Prev.Next = node.Next
			node.Next.Prev = node.Prev
			return res, nil
		} else {
			node = node.Next
		}
	}
	return nil, errors.New("Not Found ")
}