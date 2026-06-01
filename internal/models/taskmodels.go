package models

const (
	DataFormat = "20060102"
)

type TaskResponse struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id      string `json:"id,omitempty"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}
