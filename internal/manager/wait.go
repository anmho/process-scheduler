package manager

import (
	"fmt"
	"github.com/anmho/cs143b/project1/internal/resources"
)

func (m *Manager) updateWaitlist(resourceID, freedUnits int) error {
	resource := m.resources[resourceID]
	cur := resource.Waitlist.Waiting.Front()
	for cur != nil && resource.Available() > 0 {
		next := cur.Next()

		wu, ok := cur.Value.(*resources.WaitingUnits)
		if !ok {
			panic("invalid type assertion")
		}
		proc := m.processes[wu.Pid]

		if proc == nil {
			return nil
		}

		need := resource.Waitlist.WaitingFor(proc.Pid)
		if resource.Available() >= need {
			// update the quantity of the resource
			err := resource.Request(need)
			if err != nil {
				return fmt.Errorf("requesting units: %w", err)
			}
			// update held by process
			proc.AddResource(resourceID, need)
			err = m.ready.Add(wu.Pid, proc.Priority)
			if err != nil {
				return fmt.Errorf("could not add pid %d to Ready list: %w", wu.Pid, err)
			}
			resource.Waitlist.Remove(wu.Pid)
		} else { // why stop?
			break
		}
		cur = next
	}
	return nil
}
