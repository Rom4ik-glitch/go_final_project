package main

import (
	"log"
	"net/http"
	"os"
)

func ServerStart(port, dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	if install {
		CreateDB(dbFile)
	}

	dbConn, err := OpenDB(dbFile)
	if err != nil {
		return err
	}

	db = dbConn
	defer dbConn.Close() ///
	defer db.Close()     ///

	log.Printf("Server started at port %s", port)

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)

	http.HandleFunc("/api/tasks", getTasksHandler)

	http.HandleFunc("/api/nextdate", nextDateHandler)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
