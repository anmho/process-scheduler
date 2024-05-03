package process

import (
	"container/list"
	"errors"
)

type State int

const (
	ReadyState   State = 0
	WaitingState State = 1
)

type PCB struct {
	Pid       int
	Priority  int
	State     State
	Parent    int
	Children  *list.List  // double linked list of indexes
	Resources map[int]int // double linked list of indexes
}

func New(parent int, pid int, priority int) *PCB {
	return &PCB{
		Pid:       pid,
		Priority:  priority,
		State:     ReadyState,
		Parent:    parent,
		Children:  list.New(),
		Resources: make(map[int]int),
	}
}

func (p *PCB) AddChild(pid int) {
	p.Children.PushBack(pid)
}

func (p *PCB) AddResource(resourceID int, units int) {
	p.Resources[resourceID] += units
}

func (p *PCB) HoldingResource(resourceID int) int {
	return p.Resources[resourceID]
}

func (p *PCB) ReleaseResource(resourceID int, units int) error {
	if p.Resources[resourceID] < units {
		return errors.New("not enough Resources to release")
	}

	p.Resources[resourceID] -= units
	return nil
}
