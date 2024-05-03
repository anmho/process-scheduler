package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReadyList(t *testing.T) {
	assert := assert.New(t)

	rl := NewReadyList(3)
	assert.Equal(3, rl.levels)
}

func TestAdd(t *testing.T) {
	assert := assert.New(t)
	var pid int
	var ok bool
	var err error
	rl := NewReadyList(3)

	// PushBack first process
	err = rl.Add(0, 0)

	assert.Nil(err)
	assert.Equal(rl.Ready[0].Len(), 1)
	pid, ok = rl.Ready[0].Front().Value.(int)
	assert.True(ok)
	assert.Equal(pid, 0)

	// PushBack a second process
	err = rl.Add(1, 1)

	assert.Nil(err)
	assert.Equal(rl.Ready[1].Len(), 1)
	pid, ok = rl.Ready[1].Front().Value.(int)
	assert.True(ok)
	assert.Equal(pid, 1)

	pid, err = rl.Running()
	assert.Nil(err)
	assert.Equal(1, pid)

}

func TestGetRunning(t *testing.T) {
	var pid int
	var err error
	var assert = assert.New(t)
	rl := NewReadyList(3)

	// pid 0
	rl.Ready[0].PushBack(0)
	pid, err = rl.Running()
	assert.Nil(err)
	assert.Equal(pid, 0)

	rl.Ready[1].PushBack(1)
	pid, err = rl.Running()

	// pid 1
	assert.Nil(err)
	assert.Equal(pid, 1)

	// pid 2
	rl.Ready[1].PushBack(2)
	pid, err = rl.Running()

	assert.Nil(err)
	assert.Equal(pid, 1)

	// pid 3
	rl.Ready[1].PushBack(3)
	pid, err = rl.Running()

	assert.Nil(err)
	assert.Equal(pid, 1)

	// pid 3
	rl.Ready[1].PushBack(3)
	pid, err = rl.Running()

	assert.Nil(err)
	assert.Equal(pid, 1)
}
