package main

import (
	"log"
	"os"

	"github.com/St-Ivanov/todo-list-project/internal/api"
	"github.com/St-Ivanov/todo-list-project/internal/db"
	"github.com/St-Ivanov/todo-list-project/internal/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.DB.Close()

	err = os.Mkdir("./log", 0755)
	if err != nil && !os.IsExist(err) {
		log.Println(err)
		return
	}

	file, err := os.OpenFile("./log/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	loger := log.New(file, "", 0)

	serv := server.NewHttpServer(loger)

	api.Pass = os.Getenv("TODO_PASSWORD")

	err = serv.Serv.ListenAndServe()
	if err != nil {
		serv.Loger.Println(err)
		return
	}
}
