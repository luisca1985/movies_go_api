package main

import "fmt"

type taskList struct {
	tasks []*task
}

func (t *taskList) addToList(tl *task) {
	t.tasks = append(t.tasks, tl)
}

func (t *taskList) removeFromList(index int) {
	t.tasks = append(t.tasks[:index], t.tasks[index+1:]...)
}

func (t *taskList) printList() {
	for _, task := range t.tasks {
		fmt.Println("Name", task.name)
		fmt.Println("Description", task.description)
	}

}

func (t *taskList) printListCompleted() {
	for _, task := range t.tasks {
		if task.completed == true {
			fmt.Println("Name", task.name)
			fmt.Println("Description", task.description)

		}
	}

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
	list.addToList(t3)
	list.printList()
	list.tasks[0].markAsCompleted()
	fmt.Println("Tasks completed")
	list.printListCompleted()

	mapTasks := make(map[string]*taskList)

	mapTasks["Nestor"] = list

	t4 := &task{
		name:        "Complete Java course",
		description: "Complete Java course this week",
	}
	t5 := &task{
		name:        "Complete C# course",
		description: "Complete C# course this week",
	}

	list2 := &taskList{
		tasks: []*task{
			t4, t5,
		},
	}

	mapTasks["Ricardo"] = list2

	fmt.Println("Tareas de Nestor")
	mapTasks["Nestor"].printList()

	fmt.Println("Tareas de Ricardo")
	mapTasks["Ricardo"].printList()

}
