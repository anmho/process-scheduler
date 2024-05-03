package manager

import "fmt"

func (m *Manager) Scheduler() {
	pid, err := m.ready.Running()
	if err != nil {
		fmt.Printf("%d ", -1)
		return
	}

	fmt.Printf("%d ", pid)
}
