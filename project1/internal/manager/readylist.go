package manager

import "container/list"



type ReadyList struct {
	ready []*list.List
	levels int 
}

func NewReadyList(levels int) *ReadyList {
	processes := make([]*list.List, levels)
	rl := &ReadyList{
		ready: processes,
		levels: levels,
	}
	return rl
}


// Returns the pid of the running process (the highest priority ready process)
func (rl *ReadyList) GetRunning() int {
	// get the highest priority item in the readylist
	// basically get the first non-null head
	// within each list its fifo

	var pid int = -1
	for i := rl.levels; i >= 0; i-- {
		if rl.ready[i].Front() != nil {
			pid = rl.ready[i].Front().Value.(int)
			break
		}
	}

	return pid
}


// Add a process id to the ready list which has priority p
func (rl *ReadyList) Add(pid int, priority int) {
	rl.ready[priority].PushBack(pid)
}