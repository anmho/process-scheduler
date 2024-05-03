package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/anmho/cs143b/project1/internal/manager"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var m *manager.Manager = nil

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) > 0 {
			args := strings.Split(line, " ")
			err := execCmd(args, &m)
			if err != nil {
				log.Printf("cmd %s error: %w\n", args[0], err)
				fmt.Printf("%d \n", -1)
			}
		}
	}
}

func execCmd(args []string, m **manager.Manager) error {
	if m == nil {
		return errors.New("m is nil")
	}

	//fmt.Println(args)

	switch args[0] {
	case "in":
		// in <n> <u0> <u1> <u2> <u3>
		// Initialize with n priority levels in readyList
		// inventory of u0 for resources 0
		// inventory of u1 for resources 1
		// inventory of u2 for resources 2
		// inventory of u3 for resources 3
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
		if err != nil || u0 <= 0 {
			return errors.New("invalid value for u0")
		}
		u1, err := strconv.Atoi(args[3])
		if err != nil || u1 <= 0 {
			return errors.New("invalid value for u1")
		}
		u2, err := strconv.Atoi(args[4])
		if err != nil || u2 <= 0 {
			return errors.New("invalid value for u2")
		}
		u3, err := strconv.Atoi(args[5])
		if err != nil || u3 <= 0 {
			return errors.New("invalid value for u3")
		}

		if *m != nil {
			fmt.Println()
		}
		*m = manager.New(n, u0, u1, u2, u3)
	case "id":
		// id
		// Initialize with default values
		// Equivalent to in 3 1 1 2 3
		// num levels = 3
		// u0 = 1, u1 = 1, u2 =  2, u3 = 3

		if *m != nil {
			fmt.Println()
		}
		*m = manager.NewDefault()
	case "cr":
		// cr <p>
		// Create child process for running process at priority level p
		if len(args) != 2 {
			return errors.New("invalid num args for cr <p>")
		}

		priority, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("invalud value for p for cr <p>")
		}
		err = (*m).Create(priority)
		if err != nil {
			return err
		}
	case "de":
		var err error
		// de <i>
		// destroy process i and all of its descendants

		if len(args) != 2 {
			return errors.New("not enough args for de")
		}
		pid, err := strconv.Atoi(args[1])

		if err != nil {
			return errors.New("invalid value for de <i>")
		}

		err = (*m).Destroy(pid)
		if err != nil {
			return fmt.Errorf("could not destroy pid %d: %w", pid, err)
		}
	case "rq":
		// Request
		// rq <r> <k>
		// Invoke the function request(), which requests <k> units of resources <r>;
		// <r> can be 0, 1, 2, or 3.

		// if it results in a deadlock or there are not enough units then print -1
		// Invoke the function request(), which requests <k> units of resources <r>;
		// <r> can be 0, 1, 2, or 3.
		if len(args) != 3 {
			return errors.New("invalid num args for rq")
		}
		resourceID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("invalid value for r")
		}
		units, err := strconv.Atoi(args[2])
		if err != nil {
			return errors.New("invalid value for k")
		}
		err = (*m).Request(resourceID, units)
		if err != nil {
			return fmt.Errorf("could not request %d units of resources %d %w", resourceID, units, err)
		}
	case "rl":
		// Release
		// rl <r> <k>
		// Invoke the function release(), which release the resources <r>;
		// <r> can be 0, 1, 2, or 3; <k> is the number of units to be released
		if len(args) != 3 {
			return errors.New("invalid num args for rl <r> <k>")
		}

		resourceID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("invalid value for <r>")
		}
		units, err := strconv.Atoi(args[2])
		if err != nil {
			return errors.New("invalid value for <k>")
		}

		err = (*m).Release(resourceID, units)
		if err != nil {
			return fmt.Errorf("could not release resources %d %w", resourceID, err)
		}
	case "to":
		// Invoke the timer interrupt
		err := (*m).Timeout()
		if err != nil {
			return fmt.Errorf("timeout error: %w", err)
		}
	case "exit":
		os.Exit(0)
	default:
		return errors.New("invalid command")
	}

	return nil
}
