package resources

import (
	"errors"
)

type RCB struct {
	available      int // number of Units available
	TotalInventory int
	Waitlist       *Waitlist
}

func New(inventory int) *RCB {
	waitlist := NewWaitlist()
	return &RCB{
		available:      inventory, // all are Free at the start
		TotalInventory: inventory,
		Waitlist:       waitlist,
	}
}

func (r *RCB) Request(units int) error {
	if units > r.available {
		return errors.New("not enough resources")
	}

	r.available -= units
	return nil
}

// Release Processes resources request from process Pid
func (r *RCB) Release(units int) error {
	if r.available+units > r.TotalInventory {
		return errors.New("not enough resources")
	}
	r.available += units

	return nil
}

func (r *RCB) Available() int {
	return r.available
}
