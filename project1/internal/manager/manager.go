package manager

import (
	"errors"
	"fmt"
	"os"

	"github.com/anmho/cs143b/project1/internal/process"
	"github.com/anmho/cs143b/project1/internal/resource"
)

const maxProcs = 16
const maxResources = 4

type Manager struct {
	ready *ReadyList
	// waiting *list.List
	processes *[maxProcs]*process.PCB
	resources *[maxResources]*resource.RCB
	created   int
}

func New(levels int, u0 int, u1 int, u2 int, u3 int) *Manager {
	readyList := NewReadyList(levels)

	processes := new([maxProcs]*process.PCB)
	resources := new([maxResources]*resource.RCB)

	m := &Manager{
		ready: readyList,
		// waiting: wait,
		processes: processes,
		resources: resources,
		created:   0,
	}

	init := process.New(-1, 0, 0)
	m.addProcess(init)

	m.Scheduler()

	return m
}

func NewDefault() *Manager {
	return New(3, 1, 2, 2, 3)
}

func (m *Manager) addProcess(pcb *process.PCB) {
	m.processes[m.created] = pcb
	_ = m.ready.Add(pcb.Pid, pcb.Priority)
	m.created++

}

func (m *Manager) Create(priority int) {
	running := m.getRunning()
	p := process.New(running.Pid, m.created, priority)
	m.addProcess(p)
	running.Children.PushBack(p.Pid)
	m.Scheduler()
}

func (m *Manager) Destroy(pid int) error {
	p := m.processes[pid]

	for cur := p.Children.Front(); cur != nil; cur = cur.Next() {
		childPid := cur.Value.(int)
		err := m.Destroy(childPid)
		if err != nil {
			return fmt.Errorf("could not destroy pid %d: %w", childPid, err)
		}
	}
	//for cur := p.Resources
	cur := p.Resources.Front()
	for cur != nil {
		held, ok := cur.Value.(process.HeldResource)
		if !ok {
			panic("invalid type assertion process resource")
		}
		err := p.ReleaseResource(held.ResourceID, held.Units)
		if err != nil {
			return errors.New("could not release resource")
		}

		cur = cur.Next()
	}
	//p.Resources
	m.processes[pid] = nil
	m.Scheduler()
	return nil
}
func (m *Manager) Request(resource int, units int) error {
	err := m.resources[resource].Request(units)
	if err != nil {
		return fmt.Errorf("could not request resource %d: %w", resource, err)
	}

	m.Scheduler()
	return nil
}

func (m *Manager) Release(resource int, units int) {
	m.Scheduler()
}

func (m *Manager) Timeout() error {

	// move it from the front of the readylist to the back
	running := m.getRunning()
	//log.Println("removing", running.Pid)
	err := m.ready.TimerInterrupt(running.Pid, running.Priority)
	if err != nil {
		return fmt.Errorf("manager timeout error: %w", err)
	}

	//m.ready.Print()
	m.ready.Add(running.Pid, running.Priority)

	m.Scheduler()
	return nil
}

func (m *Manager) Scheduler() {
	pid, err := m.ready.Running()
	if err != nil {
		fmt.Printf("%d ", -1)
		return
	}

	fmt.Fprintf(os.Stdout, "%d ", pid)
}

func (m *Manager) getRunning() *process.PCB {
	pid, _ := m.ready.Running()

	return m.processes[pid]
}
