package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/St-Ivanov/todo-list-project/internal/db"
	"github.com/St-Ivanov/todo-list-project/internal/models"
)

const (
	limitTasksOnPage = 50
)

var (
	errNotFindID        = errors.New("ID not find")
	errIdIsEmpty        = errors.New("ID is empty.")
	errIncIdFormat      = errors.New("incorrect id format")
	errMethodNotAllowed = errors.New("Methon not allowed")
)

// Distribution handler
func handlerTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlerAddTask(w, r)
	case http.MethodGet:
		handlerGetTask(w, r)
	case http.MethodPut:
		handlerUpdateTask(w, r)
	case http.MethodDelete:
		handlerDeleteTask(w, r)
	default:
		http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
	}
}

// Task addition handler
func handlerAddTask(w http.ResponseWriter, r *http.Request) {
	task, err := readJson(r)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	err = updateDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	last, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(last, 10)})
}

// Handler for receiving all tasks
func handlerGetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	search := r.FormValue("search")

	data, err := db.GetTasks(limitTasksOnPage, search)
	if err != nil {
		writeJson(w, map[string]string{"error": "db error."})
		return
	}

	writeJson(w, data)
}

// Task update handler
func handlerUpdateTask(w http.ResponseWriter, r *http.Request) {
	task, err := readJson(r)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	err = updateDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}

// Task deletion handler
func handlerDeleteTask(w http.ResponseWriter, r *http.Request) {
	idString := r.FormValue("id")
	if idString == "undefined" {
		writeJson(w, map[string]string{"error": errIdIsEmpty.Error()})
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": errIncIdFormat.Error()})
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}

// A handler for marking a task as completed
func handlerDoneTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	idString := r.FormValue("id")
	if idString == "undefined" {
		writeJson(w, map[string]string{"error": errIdIsEmpty.Error()})
		return
	}
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": errIncIdFormat.Error()})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": errNotFindID.Error()})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		now := time.Now().UTC().Truncate(24 * time.Hour)

		data, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}

		err = db.UpdateDoneTask(id, data)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJson(w, map[string]string{})
}

// Handler for receiving a single task by ID
func handlerGetTask(w http.ResponseWriter, r *http.Request) {
	idString := r.FormValue("id")
	if idString == "undefined" {
		writeJson(w, map[string]string{"error": errIdIsEmpty.Error()})
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": errIncIdFormat.Error()})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": errNotFindID.Error()})
		return
	}
	writeJson(w, task)
}

// A handler for updating the date for a task with a specified rule for repetition
func updateDate(task *models.Task) error {
	var (
		date time.Time
		err  error
	)
	if task.Date != "" {
		date, err = time.Parse(models.DataFormat, task.Date)
		if err != nil {
			return fmt.Errorf("the date is presented in the wrong format.")
		}

		now := time.Now().UTC().Truncate(24 * time.Hour)

		var newDate string

		if task.Repeat != "" {
			newDate, err = nextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("the repetition rule is in the wrong format.")
			}
		} else {
			newDate = now.Format(models.DataFormat)
		}

		if date.Before(now) {
			task.Date = newDate
		}
	} else {
		task.Date = time.Now().UTC().Format(models.DataFormat)
	}
	return nil
}

// Function for serializing to json and sending
func writeJson(w http.ResponseWriter, data any) {
	js, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// A function for deserialize json tasks
func readJson(r *http.Request) (models.Task, error) {
	task := models.Task{}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return models.Task{}, err
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		return models.Task{}, fmt.Errorf("error deserializing JSON.")
	}

	if task.Title == "" {
		return models.Task{}, fmt.Errorf("the issue title is not specified.")
	}

	return task, nil
}
