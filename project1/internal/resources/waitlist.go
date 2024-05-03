package resources

import "container/list"

type Waitlist struct {
	Waiting *list.List
}

type WaitingUnits struct {
	Pid   int
	Units int
}

func NewWaitlist() *Waitlist {
	return &Waitlist{Waiting: list.New()}
}

func (wl *Waitlist) PushBack(pid, units int) {
	wu := &WaitingUnits{
		Pid:   pid,
		Units: units,
	}
	wl.Waiting.PushBack(wu)
}

// WaintingFor returns the number of Units the process is waiting for
func (wl *Waitlist) WaitingFor(pid int) int {
	cur := wl.Waiting.Front()
	for cur != nil {
		wu, ok := cur.Value.(*WaitingUnits)
		if !ok {
			panic("invalid type assert for waitlist Pid")
		}
		if wu.Pid == pid {
			return wu.Units
		}

		cur = cur.Next()
	}
	return 0
}

func (wl *Waitlist) Remove(pid int) {
	cur := wl.Waiting.Front()
	for cur != nil {
		wu, ok := cur.Value.(*WaitingUnits)
		if !ok {
			panic("type assertion failed")
		}
		if wu.Pid == pid {
			wl.Waiting.Remove(cur)
		}

		cur = cur.Next()
	}
}
