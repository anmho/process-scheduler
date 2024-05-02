package resource

import "container/list"

type Waitlist struct {
	Waiting *list.List
}

func NewWaitlist() *Waitlist {
	return &Waitlist{Waiting: list.New()}
}

func (wl *Waitlist) PushBack(pid int) {
	wl.Waiting.PushBack(pid)
}

func (wl *Waitlist) Remove(pid int) {
	cur := wl.Waiting.Front()
	for cur != nil {
		curPid, ok := cur.Value.(int)
		if !ok {
			panic("type assertion failed")
		}
		if curPid == pid {
			wl.Waiting.Remove(cur)
		}

		cur = cur.Next()
	}
}
