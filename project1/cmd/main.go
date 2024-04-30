package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/anmho/cs143a/project1/internal/manager"
)



func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var m *manager.Manager


	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Split(line, " ")
		err := execCmd(args, m)
		if err != nil {
			fmt.Printf("%d ", -1)
		}
	}
}

func execCmd(args []string, m *manager.Manager) error {
	switch args[0] {
	case "in":
		// in <n> <u0> <u1> <u2> <u3>
		// Initialize with n priority levels in readyList
		// inventory of u0 for resource 0
		// inventory of u1 for resource 1
		// inventory of u2 for resource 2
		// inventory of u3 for resource 3
		if len(args) != 6 {
			return errors.New("invalid number of arguments")
		}

		// Priority levels
		n, err := strconv.Atoi(args[1])
		if err != nil || n <= 0 {
			return errors.New("invalid value for n")
		}
		// Resource inventories
		u0, err := strconv.Atoi(args[2])
		if err !=  nil || u0 <= 0 {
			return errors.New("invalid value for u0")
		}
		u1, err := strconv.Atoi(args[3])
		if err !=  nil || u1 <= 0 {
			return errors.New("invalid value for u1")
		}
		u2, err := strconv.Atoi(args[4])
		if err !=  nil || u2 <= 0 {
			return errors.New("invalid value for u2")
		}
		u3, err := strconv.Atoi(args[5])
		if err !=  nil || u3 <= 0 {
			return errors.New("invalid value for u3")
		}

		if m == nil {
			fmt.Println()
		}
		m = manager.New(n, u0, u1, u2, u3)
	case "id":
		// id 
		// Initialize with default values
		// Equivalent to in 3 1 1 2 3
		// num levels = 3
		// u0 = 1, u1 = 1, u2 =  2, u3 = 3
		if m == nil {
			fmt.Println()
		}
		m = manager.New(3, 1, 2, 2, 3)
	case "cr":
		// cr <p>
		// Create child process for running process at priority level p
		m.Create()
	case "de":
		// de <i>
		// destroy process i and all of its descendants
		if len(args) != 2 {
			return errors.New("not enough args for de")
		}
		pid, err := strconv.Atoi(args[1])

		if err != nil {
			return errors.New("invalid value for de <i>")
		}

		m.Destroy(pid)
	case "rq":
		// Request
		// rq <r> <k>
		// Invoke the function request(), which requests <k> units of resource <r>; 
		// <r> can be 0, 1, 2, or 3.

		// if it results in a deadlock or there are not enough units then print -1
		// Invoke the function request(), which requests <k> units of resource <r>; 
		// <r> can be 0, 1, 2, or 3.
		if len(args) != 3 {
			return errors.New("invalid num args for rq")
		}
		resourceNum, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("invalid value for r")
		}
		units, err := strconv.Atoi(args[2])
		if err != nil {
			return errors.New("invalid value for k")
		}
		m.Request(resourceNum, units)
	case "rl":
		// Release
		// rl <r> <k>
		// Invoke the function release(), which release the resource <r>;
		// <r> can be 0, 1, 2, or 3; <k> is the number of units to be released
		if len(args) != 3 {
			return errors.New("invalid num args for rl <r> <k>")
		}

		resourceNum, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("invalid value for <r>")
		}
		units, err := strconv.Atoi(args[2])
		if err != nil {
			return errors.New("invalid value for <k>")
		}

		m.Release(resourceNum, units)
	case "to":
		// Invoke the timer interrupt
		m.Timeout()
	case "exit":
		os.Exit(0)
	default:
		return errors.New("invalid command")
	}

	m.Scheduler()
	return nil
}