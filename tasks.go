package main

import (
	"fmt"
)

type taskList struct {
	tasks []*task
}

func (t *taskList) addToList(tl *task) {
	t.tasks = append(t.tasks, tl)
}

func (t *taskList) removeFromList(index int) {
	t.tasks = append(t.tasks[:index], t.tasks[index+1:]...)
}

type task struct {
	name        string
	description string
	completed   bool
}

func (t *task) markAsCompleted() {
	t.completed = true
}

func (t *task) updateDescription(description string) {
	t.description = description
}

func (t *task) updateName(name string) {
	t.name = name
}

func main() {
	t1 := &task{
		name:        "Complete Go course",
		description: "Complete Go course this week",
	}
	t2 := &task{
		name:        "Complete Python course",
		description: "Complete Python course this week",
	}
	t3 := &task{
		name:        "Complete NodeJS course",
		description: "Complete NodeJs course this week",
	}

	list := &taskList{
		tasks: []*task{
			t1, t2,
		},
	}
	fmt.Println(list.tasks[0])
	list.addToList(t3)
	fmt.Println(len(list.tasks))
	list.removeFromList(1)
	fmt.Println(len(list.tasks))
}
