package manager

import (
	"github.com/anmho/cs143a/project1/internal/process"
	"github.com/anmho/cs143a/project1/internal/resource"
)

const max_procs = 16
const max_resources = 4

type Manager struct {
	ready *ReadyList
	// waiting *list.List
	processes *[max_procs]*process.PCB
	resources *[max_resources]*resource.RCB
	created int
}

func New(levels int, u0 int, u1 int, u2 int, u3 int) *Manager {
	readyList := NewReadyList(levels)
	
	processes := new([max_procs]*process.PCB)
	resources := new([max_resources]*resource.RCB)


	m := &Manager{
		ready: readyList,
		// waiting: wait,
		processes: processes,
		resources: resources,
		created: 0,
	}
	return m
}

func (m *Manager) Create() {
	running_idx := m.ready.GetRunning()
	running := m.processes[running_idx]

	pcb := process.New(running.Pid, m.created)
	m.created++

	m.processes[pcb.Pid] = pcb
}


func (m *Manager) Destroy(pid int) {
	running_idx := m.ready.GetRunning()

	running := m.processes[running_idx]

	// destroy running process and all children
	// should probably be tail

	for child := running.Children.Front(); child != nil; child = child.Next() {
		child_idx := child.Value.(int)
		m.Destroy(child_idx)
	}

	// free, then set to nil

	m.processes[running_idx] = nil
}
func (m *Manager) Request(resource int, units int) {

}

func (m *Manager) Release(resource int, units int) {

}

func (m *Manager) Timeout() {
	// 

}

func (m *Manager) Scheduler() {


}