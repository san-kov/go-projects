package main

import (
	"fmt"
	"os"
)

type Task struct {
	ID   int
	Text string
	Done bool
}

func main() {

	var tasks []Task

	if len(os.Args) < 2 {
		fmt.Println("Usage: task-manager <command> [args]")
		return
	}

}

func addTask(tasks []Task, text string) []Task {
	id := len(tasks) + 1
	newTask := Task{ID: id, Text: text, Done: false}
	tasks = append(tasks, newTask)

	return tasks
}

func completeTask(tasks []Task, id int) []Task {

	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		if task.ID == id {

			tasks[i].Done = true
			return tasks
		}
	}

	fmt.Println("Task not found")

	return tasks
}

func printTasks(tasks []Task) {
	for i := 0; i < len(tasks); i++ {
		var check string
		if tasks[i].Done {
			check = "x"
		} else {
			check = " "
		}
		fmt.Printf("[%s] %d %s\n", check, tasks[i].ID, tasks[i].Text)
	}
}
