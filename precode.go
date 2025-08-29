package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

type TaskReturnWithMesssage struct {
	Task
	Message string `json:"message"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
func getTasks(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена!", http.StatusInternalServerError)
		return
	}
	data, jsonErr := json.Marshal(task)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	_, writeErr := w.Write(data)
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}
func createTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer r.Body.Close()

    var task Task
    err = json.Unmarshal(body, &task)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	tasks[task.ID] = task
	message := TaskReturnWithMesssage{
		Task:task,
		Message: "Задача успещно добавлена!",
	}

	data, jsonErr := json.Marshal(message)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	_, writeErr := w.Write(data)
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена!", http.StatusInternalServerError)
		return
	}

	delete(tasks, id)
	message := TaskReturnWithMesssage{
		Task:task,
		Message: "Задача успещно удалена!",
	}

	data, jsonErr := json.Marshal(message)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	_, writeErr := w.Write(data)
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Get("/tasks/{id}", getTask)
	r.Post("/tasks", createTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
