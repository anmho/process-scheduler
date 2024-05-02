package resource

import (
	"errors"
)

type RCB struct {
	Available      int // number of units available
	TotalInventory int
	Waitlist       *Waitlist
}

func New(inventory int) *RCB {
	waitlist := NewWaitlist()
	return &RCB{
		Available:      inventory, // all are Free at the start
		TotalInventory: inventory,
		Waitlist:       waitlist,
	}
}

func (r *RCB) Request(units int) error {
	if units > r.Available {
		return errors.New("not enough resources")
	}

	r.Available -= units
	return nil
}

// Release Processes resource request from process pid
func (r *RCB) Release(units int) error {
	if r.Available+units > r.TotalInventory {
		return errors.New("not enough resources")
	}
	r.Available += units

	return nil
}
