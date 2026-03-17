package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Task struct {
	ID   int
	Text string
	Done bool
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: task-manager <command> [args]")
		return
	}

	tasks, err := loadTasks()

	if err != nil {
		fmt.Println("Failed to load tasks")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-manager add <text>")
			return
		}

		text := os.Args[2]
		tasks = addTask(tasks, text)

		if err := saveTasks(tasks); err != nil {
			fmt.Println("Failed to save tasks")
			return
		}

		fmt.Printf("Added task %d: %s\n", len(tasks), text)
	case "list":
		printTasks(tasks)
	default:
		fmt.Printf("Unknown command: %s\n", command)
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

func saveTasks(tasks []Task) error {
	data, err := json.Marshal(tasks)

	if err != nil {
		return err
	}

	os.WriteFile("tasks.json", data, 0644)

	return nil
}

func loadTasks() ([]Task, error) {
	var loaded []Task
	data, err := os.ReadFile("tasks.json")

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}

		return nil, err
	}

	json.Unmarshal(data, &loaded)

	return loaded, nil
}
