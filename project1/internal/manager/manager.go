package manager

import (
	"errors"
	"fmt"
	"github.com/anmho/cs143b/project1/internal/process"
	"github.com/anmho/cs143b/project1/internal/resources"
)

const maxProcesses = 16
const maxResources = 4

type Manager struct {
	ready     *ReadyList
	levels    int
	processes *[maxProcesses]*process.PCB
	resources *[maxResources]*resources.RCB
	created   int
}

func (m *Manager) initResources(u0, u1, u2, u3 int) {
	m.resources = new([maxResources]*resources.RCB)
	m.resources[0] = resources.New(u0)
	m.resources[1] = resources.New(u1)
	m.resources[2] = resources.New(u2)
	m.resources[3] = resources.New(u3)
}
func New(levels, u0, u1, u2, u3 int) *Manager {
	readyList := NewReadyList(levels)

	processes := new([maxProcesses]*process.PCB)
	resources := new([maxResources]*resources.RCB)

	m := &Manager{
		ready:     readyList,
		levels:    levels,
		processes: processes,
		resources: resources,
		created:   0,
	}

	m.initResources(u0, u1, u2, u3)

	init := process.New(-1, 0, 0)
	m.addProcess(init)
	m.Scheduler()

	return m
}

func NewDefault() *Manager {
	return New(3, 1, 1, 2, 3)
}

func (m *Manager) addProcess(pcb *process.PCB) {
	m.processes[m.created] = pcb
	_ = m.ready.Add(pcb.Pid, pcb.Priority)
	m.created++

}

func (m *Manager) Create(priority int) error {
	if !m.isValidPriority(priority) {
		return fmt.Errorf("invalid priority: %d", priority)
	}
	running := m.getRunning()
	p := process.New(running.Pid, m.created, priority)
	m.addProcess(p)
	running.Children.PushBack(p.Pid)
	m.Scheduler()
	return nil
}

// Returns the child pid if it was found as a descendent of pid
func (m *Manager) isDescendent(pid, targetPid int) bool {
	cur := m.processes[pid]
	if cur == nil {
		return false
	}
	if pid == targetPid {
		return true
	}
	for child := cur.Children.Front(); child != nil; child = child.Next() {
		childPid, ok := child.Value.(int)
		if !ok {
			panic("invalid child pid type assertion")
		}

		if m.isDescendent(childPid, targetPid) {
			return true
		}
	}
	return false
}

func (m *Manager) destroy(pid int) error {
	proc := m.processes[pid]
	if proc == nil {
		return nil
	}
	//log.Printf("removing children of %d\n", pid)
	for child := proc.Children.Front(); child != nil; child = child.Next() {
		childPid, ok := child.Value.(int)
		if !ok {
			panic("invalid child pid type assertion")
		}
		//log.Printf("%d is a child of %d", childPid, proc.Pid)
		err := m.destroy(childPid)
		if err != nil {
			return fmt.Errorf("destroying %d: %w", proc.Pid, err)
		}
	}

	err := m.ready.Remove(pid, proc.Priority)
	if err != nil {
		return fmt.Errorf("could not remove pid %d from Ready list: %w", pid, err)
	}

	for resourceID, hold := range proc.Resources {
		err := m.release(proc.Pid, resourceID, hold)
		if err != nil {
			return err
		}
	}

	m.processes[pid] = nil

	return nil
}

func (m *Manager) isValidPriority(priority int) bool {
	return priority >= 0 && priority <= m.levels-1
}

func isValidPID(pid int) bool {
	return pid >= 0 && pid < maxProcesses
}
func (m *Manager) Destroy(pid int) error {
	running := m.getRunning()
	if running.Pid == 0 {
		return errors.New("init process cannot destroy itself")
	}
	if !isValidPID(pid) {
		return fmt.Errorf("invalid pid %d", pid)
	}
	if !m.isDescendent(running.Pid, pid) {
		return fmt.Errorf("process %d is not a descdent of running process %d", pid, running.Pid)
	}

	err := m.destroy(pid)
	if err != nil {
		return err
	}

	m.Scheduler()
	return nil
}

func (m *Manager) getRunning() *process.PCB {
	pid, _ := m.ready.Running()

	return m.processes[pid]
}
