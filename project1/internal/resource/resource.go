package resource

import (
	"container/list"
	"errors"
)

type RCB struct {
	free           bool
	state          int // number of units available
	totalInventory int
	waitlist       *list.List
}

func New(inventory int) *RCB {
	waitlist := list.New()
	return &RCB{
		free:           true,
		state:          inventory, // all are free at the start
		totalInventory: inventory,
		waitlist:       waitlist,
	}
}

func (r *RCB) Request(units int) error {
	if units > r.state {
		return errors.New("not enough resources")
	}

	r.state -= units
	return nil
}

func (r *RCB) Release(units int) {
	r.state
}
