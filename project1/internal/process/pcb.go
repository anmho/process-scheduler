package process

import "container/list"

type ProcessState int

const (
	ReadyState ProcessState = 0
	WaitingState ProcessState = 1
)

type PCB struct {
	Pid int
	State ProcessState
	Parent int
	Children *list.List // double linked list of indexes
	Resources *list.List // double linked list of indexes
}


func New(parent int, pid int) *PCB {
	return &PCB{
		Pid: pid,
		State:     ReadyState,
		Parent:    parent,
		Children:  list.New(),
		Resources: list.New(),
	}
}