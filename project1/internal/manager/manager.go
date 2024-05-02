package manager

import (
	"fmt"
	"github.com/anmho/cs143b/project1/internal/process"
	"github.com/anmho/cs143b/project1/internal/resource"
)

const maxProcesses = 16
const maxResources = 4

type Manager struct {
	ready     *ReadyList
	processes *[maxProcesses]*process.PCB
	resources *[maxResources]*resource.RCB
	created   int
}

func (m *Manager) initResources(u0, u1, u2, u3 int) {
	m.resources = new([maxResources]*resource.RCB)
	m.resources[0] = resource.New(u0)
	m.resources[1] = resource.New(u1)
	m.resources[2] = resource.New(u2)
	m.resources[3] = resource.New(u3)
}
func New(levels, u0, u1, u2, u3 int) *Manager {
	readyList := NewReadyList(levels)

	processes := new([maxProcesses]*process.PCB)
	resources := new([maxResources]*resource.RCB)

	m := &Manager{
		ready:     readyList,
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

	// Destroy all of its children first
	for cur := p.Children.Front(); cur != nil; cur = cur.Next() {
		childPid := cur.Value.(int)
		err := m.Destroy(childPid)
		if err != nil {
			return fmt.Errorf("could not destroy pid %d: %w", childPid, err)
		}
	}

	// Release the resources of the process
	for resourceID := range m.resources {
		//resourceID
		units := p.HoldingResource(resourceID)
		err := p.ReleaseResource(resourceID, units)
		if err != nil {
			return fmt.Errorf("could not release resources of pid %d: %d %w", pid, resourceID, err)
		}
	}

	m.processes[pid] = nil
	m.Scheduler()
	return nil
}

func isValidResource(resourceID int) bool {
	return resourceID >= 0 && resourceID <= maxResources
}
func (m *Manager) Request(resourceID int, units int) error {
	if !isValidResource(resourceID) {
		return fmt.Errorf("invalid resource id %d", resourceID)
	}
	running := m.getRunning()
	// check if we can request it before requesting

	resource := m.resources[resourceID]

	if resource.Available >= units {
		err := resource.Request(units)
		if err != nil {
			return fmt.Errorf("requesting units: %w", err)
		}
		running.AddResource(resourceID, units)
	} else {
		running.State = process.WaitingState
		// remove it from the ready list and put it on the waiting list
		err := m.ready.Remove(running.Pid, running.Priority)
		if err != nil {
			return fmt.Errorf("removing %d from ready list: %w", running.Pid, err)
		}
		resource.Waitlist.PushBack(running.Pid)
	}

	m.Scheduler()
	return nil
}

func (m *Manager) Release(resourceID int, units int) error {
	if !isValidResource(resourceID) {
		return fmt.Errorf("invalid resource id %d", resourceID)
	}
	var err error
	// release the resources from the process
	running := m.getRunning()
	resource := m.resources[resourceID]

	err = running.ReleaseResource(resourceID, units)
	if err != nil {
		return fmt.Errorf("could not release resources of pid %d: %d %w", resourceID, units, err)
	}

	// add them back to the resourceID pool
	err = resource.Release(units)
	if err != nil {
		return fmt.Errorf("could not release %d units of resources resourceID %d: %w", units, resourceID, err)
	}

	//cur := resource.Waitlist.Waiting.Front()
	//for resource.Waitlist.Waiting.Front() != nil && resource.Available > 0 {
	//
	//	pid, ok := cur.Value.(int)
	//	if !ok {
	//		panic("invalid type assertion")
	//	}
	//
	//
	//	cur = cur.Next()
	//}

	m.Scheduler()
	return nil
}

func (m *Manager) Timeout() error {
	var err error

	running := m.getRunning()
	err = m.ready.Remove(running.Pid, running.Priority)
	if err != nil {
		return fmt.Errorf("manager timeout error: %w", err)
	}

	err = m.ready.Add(running.Pid, running.Priority)
	if err != nil {
		return fmt.Errorf("manager timeout error when adding running process to back of readylist: %w", err)
	}

	m.Scheduler()
	return nil
}

func (m *Manager) Scheduler() {
	pid, err := m.ready.Running()
	if err != nil {
		fmt.Printf("%d ", -1)
		return
	}

	fmt.Printf("%d ", pid)
}

func (m *Manager) getRunning() *process.PCB {
	pid, _ := m.ready.Running()

	return m.processes[pid]
}
