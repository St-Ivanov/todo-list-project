package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/St-Ivanov/todo-list-project/internal/models"
)

var (
	errDBQuery = errors.New("database query error.")
	errIncId   = errors.New("incorect ID")
)

const (
	formatDataSearch = "02.01.2006"
)

// A function that queries the database to retrieve tasks
func GetTasks(limit int, search string) (models.TaskResponse, error) {
	var data *sql.Rows

	date, err := time.Parse(formatDataSearch, search)
	if err == nil && date.Format(formatDataSearch) == search {
		dateString := date.Format(models.DataFormat)
		data, err = DB.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = :date
			ORDER BY date
			LIMIT :limit;
			`,
			sql.Named("date", dateString),
			sql.Named("limit", limit),
		)
	} else if search != "" {
		search = "%" + search + "%"

		data, err = DB.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE :search OR comment LIKE :search
			ORDER BY date
			LIMIT :limit;
		`,
			sql.Named("search", search),
			sql.Named("limit", limit),
		)
	} else {
		data, err = DB.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			LIMIT :limit;
			`,
			sql.Named("limit", limit),
		)
	}
	if err != nil {
		return models.TaskResponse{}, errDBQuery
	}

	var result models.TaskResponse
	result.Tasks = make([]models.Task, 0)

	for data.Next() {
		var temp models.Task
		err = data.Scan(&temp.Id, &temp.Date, &temp.Title, &temp.Comment, &temp.Repeat)
		if err != nil {
			return models.TaskResponse{}, errDBQuery
		}
		result.Tasks = append(result.Tasks, temp)
	}

	if err = data.Err(); err != nil {
		return models.TaskResponse{}, err
	}

	return result, nil
}

// A function that queries the database to retrieve a single task by its ID
func GetTask(id int64) (*models.Task, error) {
	var task models.Task
	err := DB.QueryRow(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = :id;
	`, sql.Named("id", id),
	).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return &task, err
	}

	return &task, nil
}

// A function that queries the database to update an existing task
func UpdateTask(task *models.Task) error {
	result, err := DB.Exec(`
		UPDATE scheduler
		SET title = :title,
			date = :date,
			comment = :comment,
			repeat = :repeat
		WHERE id = :id;
		`,
		sql.Named("id", task.Id),
		sql.Named("title", task.Title),
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errIncId
	}

	return nil
}

// A function that queries the database to add a task
func AddTask(task *models.Task) (int64, error) {
	data, err := DB.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (:date, :title, :comment, :repeat);
		`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, errDBQuery
	}
	last, err := data.LastInsertId()
	if err != nil {
		return 0, errDBQuery
	}
	return last, nil
}

// A function that queries the database to delete a task from it
func DeleteTask(id int64) error {
	ans, err := DB.Exec(`
		DELETE
		FROM scheduler
		WHERE id = :id
	`, sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := ans.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errIncId
	}

	return nil
}

// A function that queries the database to update a task with a specified recurrence rule
func UpdateDoneTask(id int64, date string) error {
	ans, err := DB.Exec(`
		UPDATE scheduler
		SET date = :date
		WHERE id = :id;
		`,
		sql.Named("id", id),
		sql.Named("date", date),
	)
	if err != nil {
		return err
	}

	count, err := ans.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errIncId
	}

	return nil
}
