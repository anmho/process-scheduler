package manager

import (
	"container/list"
	"errors"
	"fmt"
)

type ReadyList struct {
	Ready  []*list.List
	levels int
}

func NewReadyList(levels int) *ReadyList {
	processes := make([]*list.List, levels)
	for i := 0; i < levels; i++ {
		processes[i] = list.New()
	}

	rl := &ReadyList{
		Ready:  processes,
		levels: levels,
	}

	return rl
}

// Running returns the pid of the running process (the highest priority Ready process)
func (rl *ReadyList) Running() (int, error) {
	// get the highest priority item in the readylist
	// basically get the first non-null head
	// within each list its fifo

	for i := rl.levels - 1; i >= 0; i-- {
		if rl.Ready[i].Front() != nil {
			pid, ok := rl.Ready[i].Front().Value.(int)
			if !ok {
				panic("type assertion failed")
			}
			return pid, nil
		}
	}

	return -1, errors.New("no running process")
}

func (rl *ReadyList) Print() {
	for i, plist := range rl.Ready {
		fmt.Printf("Priority %d: ", i)
		for cur := plist.Front(); cur != nil; cur = cur.Next() {
			if value, ok := cur.Value.(int); !ok {
				panic("invalid type assertion")
			} else {
				fmt.Printf("%d ", value)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

// Add a process id to the Ready list which has priority p
func (rl *ReadyList) Add(pid int, priority int) error {
	if !rl.isValidPriority(priority) {
		return errors.New("invalid priority level")
	}
	rl.Ready[priority].PushBack(pid)
	return nil
}

func (rl *ReadyList) isValidPriority(priority int) bool {
	return priority >= 0 && priority < rl.levels
}

func (rl *ReadyList) PrioLen(priority int) int {
	return rl.Ready[priority].Len()
}

func (rl *ReadyList) Remove(pid int, priority int) error {
	// find and remove it from the list
	if !rl.isValidPriority(priority) {
		return errors.New("invalid priority level")
	}

	//var head = rl.Ready[priority].Front()
	//for cur := head; cur != nil; cur = cur.Next() {
	//	curPid, ok := head.Value.(int)
	//	if !ok {
	//		panic("invalid type assertion")
	//	}
	//	if curPid == pid {
	//		log.Printf("removing %v\n", cur)
	//		val := rl.Ready[priority].Remove(cur)
	//		if val == nil {
	//			return fmt.Errorf("removed invalid priority %d", curPid)
	//		}
	//
	//		return nil
	//	}
	//}
	//
	//return nil
	cur := rl.Ready[priority].Front()
	for cur != nil {
		//log.Printf("ready prio %d: %d\n", proc.Priority, cur.Value.(int))
		curPID, ok := cur.Value.(int)
		if !ok {
			panic("invalid child pid type assertion")
		}
		if curPID == pid {
			rl.Ready[priority].Remove(cur)
			return nil
		}

		cur = cur.Next()
	}
	return nil
}
