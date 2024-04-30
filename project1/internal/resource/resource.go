package resource

import "container/list"



type RCB struct {
	free bool
	state int // number of units available
	inventory int
	waitlist *list.List
}

func New(inventory int) *RCB {
	waitlist := list.New()
	return &RCB{
		free:     true,
		state: inventory, // all are free at the start
		inventory: inventory,
		waitlist: waitlist,
	}
}

func (r *RCB) Request() {
}

func (r *RCB) Release() {
}