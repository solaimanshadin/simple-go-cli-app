package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct { 
	ID int `json:"id"`
	Description string `json:"description"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func printTaskList(tasks []Task) {
	for _, task:= range tasks {
		fmt.Printf("%d - %s (%s)\n", task.ID, task.Description, task.Status)
	}
}

func printFilteredTaskList(tasks []Task, status string) {
	for _, task:= range tasks {
		if task.Status == status {
			fmt.Printf("%d - %s\n", task.ID, task.Description)

		}
	}
}

func generateTaskId() int {
	data, err := os.ReadFile("id")
	if err != nil {
		data = []byte("0")
		err := os.WriteFile("id", data, 0644)
		
		if err != nil {
			panic("Failed to update id on disk")
		}
	}
	
	lastId, err := strconv.Atoi(string(data)) 
    if err != nil {
        panic("Failed to convert ID to int")
    }

    newId := lastId + 1

	err = os.WriteFile("id", []byte(strconv.Itoa(newId)), 0644)
	if err != nil {
		panic("Failed to update id on disk")
	}
	return newId
}

func addTask(tasks *[]Task, taskTitle string)  {
	newTask := Task{}
	newTask.ID =  generateTaskId()
	newTask.Description = taskTitle
	newTask.Status = "todo"
	newTask.CreatedAt = time.Now()
	newTask.UpdatedAt = time.Now()

	*tasks = append(*tasks, newTask)
}

func getTaskIndex(tasks []Task, taskId int)  int {
	foundIndex := -1
	for i , task := range tasks {
		if(task.ID == taskId) {
			foundIndex = i
		}
	}

	return foundIndex
}

func updateTaskStatus(tasks *[]Task, id int, status string) {
	index := getTaskIndex(*tasks, id)
	if index == -1 {
		panic("No task found by this task ID")
	}
    (*tasks)[index].Status = status
}

func updateTaskDescription(tasks *[]Task, id int, description string) {
	index := getTaskIndex(*tasks, id)
	if index == -1 {
		panic("No task found by this task ID")
	}
    (*tasks)[index].Description = description
}

func deleteTask(tasks *[]Task, taskId int) {
	foundIndex := getTaskIndex(*tasks, taskId)
	
	if foundIndex == -1 {
		panic("No task found by this task ID")
	}

	*tasks = append((*tasks)[:foundIndex], (*tasks)[foundIndex+1:]...)
}

func saveOnFile(tasks []Task) {
	tasksJson, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		panic("Failed to write on file")
	}
	os.WriteFile("./data.json", tasksJson, 0777)
}

func main() {
	options := os.Args
	cmd := options[1]

	data, err := os.ReadFile("./data.json")
	if err != nil {
		data = []byte(`[]`)
		os.WriteFile("./data.json", data, 0644)
	}

	tasks := []Task{}
	
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		panic("Failed to parse json")
	}


	switch cmd { 
		case "list":
			if len(os.Args) > 2 {
				printFilteredTaskList(tasks, os.Args[2])
			} else {
				printTaskList(tasks)

			}
		case "add":
			taskTitle := os.Args[2]
			addTask(&tasks, taskTitle)
			saveOnFile(tasks)
		case "delete":
			id := os.Args[2]
			idInt, err := strconv.Atoi(id)
			if err != nil {
				panic("Failed to convert ID to staring")
			}
			deleteTask(&tasks, idInt)
			saveOnFile(tasks)
		case "mark-in-progress":
			id := os.Args[2]
			idInt, err := strconv.Atoi(id)
			if err != nil {
				panic("Failed to convert ID to staring")
			}
			updateTaskStatus(&tasks, idInt, "in-progress")
			saveOnFile(tasks)
		case "mark-done":
			id := os.Args[2]
			idInt, err := strconv.Atoi(id)
			if err != nil {
				panic("Failed to convert ID to staring")
			}
			updateTaskStatus(&tasks, idInt, "done")
			saveOnFile(tasks)
		case "update":
			id := os.Args[2]
			description := os.Args[3]
			idInt, err := strconv.Atoi(id)

			if err != nil {
				panic("Failed to convert ID to staring")
			}
			updateTaskDescription(&tasks, idInt, description)
			saveOnFile(tasks)

	}

}