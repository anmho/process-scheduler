package manager

import (
	"fmt"
	"github.com/anmho/cs143b/project1/internal/process"
	"github.com/anmho/cs143b/project1/internal/resources"
)

const maxProcesses = 16
const maxResources = 4

type Manager struct {
	ready     *ReadyList
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
	if !isValidPriority(priority) {
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

func isValidPriority(priority int) bool {
	return priority >= 0 && priority <= 2
}

func isValidPID(pid int) bool {
	return pid >= 0 && pid < maxProcesses
}
func (m *Manager) Destroy(pid int) error {
	running := m.getRunning()
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

func isValidResource(resourceID int) bool {
	return resourceID >= 0 && resourceID <= maxResources
}
func (m *Manager) Request(resourceID int, units int) error {
	running := m.getRunning()

	if !isValidResource(resourceID) {
		return fmt.Errorf("invalid resources id %d", resourceID)
	}
	// check if we can request it before requesting

	resource := m.resources[resourceID]

	if resource.Available() >= units {
		err := resource.Request(units)
		if err != nil {
			return fmt.Errorf("requesting units: %w", err)
		}
		running.AddResource(resourceID, units)
	} else {
		running.State = process.WaitingState
		// remove it from the Ready list and put it on the waiting list
		err := m.ready.Remove(running.Pid, running.Priority)
		if err != nil {
			return fmt.Errorf("removing %d from Ready list: %w", running.Pid, err)
		}
		resource.Waitlist.PushBack(running.Pid, units)
	}

	m.Scheduler()
	return nil
}

func (m *Manager) release(pid, resourceID, units int) error {
	var err error
	resource := m.resources[resourceID]
	proc := m.processes[pid]

	// Release the resource from the proc
	err = proc.ReleaseResource(resourceID, units)
	if err != nil {
		return fmt.Errorf("could not release resources of pid %d: %d %w", resourceID, units, err)
	}

	// free the units back to the resourceID pool
	err = resource.Release(units)
	if err != nil {
		return fmt.Errorf("could not release %d units of resources resourceID %d: %w", units, resourceID, err)
	}

	// Update the waitlist and remove any that can be unblocked
	err = m.updateWaitlist(resourceID, units)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) updateWaitlist(resourceID, freedUnits int) error {
	resource := m.resources[resourceID]
	cur := resource.Waitlist.Waiting.Front()
	for cur != nil && resource.Available() > 0 {

		wu, ok := cur.Value.(*resources.WaitingUnits)
		if !ok {
			panic("invalid type assertion")
		}
		proc := m.processes[wu.Pid]

		if proc == nil {
			return nil
		}

		//log.Println("proc", proc)
		need := resource.Waitlist.WaitingFor(proc.Pid)
		if freedUnits >= need {
			err := resource.Request(need)
			if err != nil {
				return fmt.Errorf("requesting units: %w", err)
			}
			proc.AddResource(resourceID, need)
			err = m.ready.Add(wu.Pid, proc.Priority)
			if err != nil {
				return fmt.Errorf("could not add pid %d to Ready list: %w", wu.Pid, err)
			}
			resource.Waitlist.Remove(wu.Pid)
		} else {
			break
		}

		cur = cur.Next()
	}
	return nil
}

func (m *Manager) Release(resourceID int, units int) error {
	if !isValidResource(resourceID) {
		return fmt.Errorf("invalid resources id %d", resourceID)
	}
	running := m.getRunning()
	err := m.release(running.Pid, resourceID, units)
	if err != nil {
		return fmt.Errorf("could not release resources of running process: %w", err)
	}

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

	//m.ready.Print()
	m.Scheduler()

	return nil
}

func (m *Manager) Scheduler() {
	pid, err := m.ready.Running()
	if err != nil {
		fmt.Printf("%d \n", -1)
		return
	}

	fmt.Printf("%d \n", pid)
}

func (m *Manager) getRunning() *process.PCB {
	pid, _ := m.ready.Running()

	return m.processes[pid]
}
