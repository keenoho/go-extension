package extension

import (
	"fmt"
	"math/rand"
	"time"
)

type TaskHandler func(action string, data any) error

type Task struct {
	Id     string
	Action string
	Data   any
}

type TaskQueue struct {
	queue       []Task
	state       int // 0: idle, 1: running
	TaskHandler TaskHandler
}

func (t *TaskQueue) Init(taskHandler TaskHandler) {
	t.TaskHandler = taskHandler
}

func (t *TaskQueue) runTask() {
	if len(t.queue) < 1 {
		return
	}
	t.state = 1
	task := t.queue[0]
	t.TaskHandler(task.Action, task.Data)
	t.queue = t.queue[1:]
	t.runTask()
}

func (t *TaskQueue) Append(action string, data any) {
	task := Task{
		Id:     fmt.Sprintf("%d-%d", rand.Intn(1e9), time.Now().UnixMilli()),
		Action: action,
		Data:   data,
	}
	t.queue = append(t.queue, task)
	t.runTask()
}
