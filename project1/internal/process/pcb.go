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
	Children  *list.List // double linked list of indexes
	Resources *list.List // double linked list of indexes
}

type HeldResource struct {
	ResourceID int
	Units      int
}

func New(parent int, pid int, priority int) *PCB {
	return &PCB{
		Pid:       pid,
		Priority:  priority,
		State:     ReadyState,
		Parent:    parent,
		Children:  list.New(),
		Resources: list.New(),
	}
}

func (p *PCB) AddChild(pid int) {
	p.Children.PushBack(pid)
}

func (p *PCB) AddResource(resourceID int, units int) {
	held := &HeldResource{resourceID, units}
	p.Resources.PushBack(held)
}

func (p *PCB) ReleaseResource(resourceID int, units int) error {

	cur := p.Resources.Front()

	for cur != nil {
		held, ok := cur.Value.(*HeldResource)
		if !ok {
			return errors.New("invalid type assertion")
		}

		if held.ResourceID == resourceID {
			if units > held.Units {
				return errors.New("attempting to release more Units than is held")
			}

			if held.Units == units {
				p.Resources.Remove(cur)
				break
			} else {
				held.Units -= units
				break
			}
		}

		cur = cur.Next()
	}
	return nil
}
