package main

import (
	"fmt"
)

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
	t := &task{
		name:        "Complete Go course",
		description: "Complete Go course this week",
	}
	fmt.Println(t)
	fmt.Printf("%+v\n", t)
	t.markAsCompleted()
	t.updateName("End the Go course.")
	t.updateDescription("End the course as soon as possible.")
	fmt.Printf("%+v\n", t)
}
