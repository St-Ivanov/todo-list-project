package api

import (
	"net/http"
)

const (
	webDir = "./web"
)

// Installing basic handlers for endpoints
func Init(h *http.ServeMux) {
	h.Handle("/", http.FileServer(http.Dir(webDir)))
	h.HandleFunc("/api/nextdate", handlerNextDate)
	h.HandleFunc("/api/task", checkAuth(handlerTask))
	h.HandleFunc("/api/tasks", checkAuth(handlerGetTasks))
	h.HandleFunc("/api/task/done", checkAuth(handlerDoneTask))
	h.HandleFunc("/api/signin", handlerAuth)
}
