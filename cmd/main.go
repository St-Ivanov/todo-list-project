package main

import (
	"log"
	"os"

	"github.com/St-Ivanov/todo-list-project/internal/db"
	"github.com/St-Ivanov/todo-list-project/internal/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir("./log", 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	file, err := os.OpenFile("./log/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	loger := log.New(file, "", 0)

	serv := server.NewHttpServer(loger)

	err = serv.Serv.ListenAndServe()
	if err != nil {
		serv.Loger.Fatal(err)
	}
}
