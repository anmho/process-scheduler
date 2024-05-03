package manager

import "fmt"

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
