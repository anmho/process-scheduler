package manager

import (
	"fmt"
	"github.com/anmho/cs143b/project1/internal/process"
)

func isValidResource(resourceID int) bool {
	return resourceID >= 0 && resourceID <= maxResources
}
func (m *Manager) Request(resourceID int, units int) error {
	running := m.getRunning()
	if units <= 0 {
		return fmt.Errorf("invalid units %d", units)
	}

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

// release releases units of resourceID from pid
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
