package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

type Config struct {
	dBFile string
	port   string
}

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

func (t *Task) SetDate() (bool, string, *Task) {
	newTask := t

	todayDateStr := time.Now().Format(dataFormat)
	if len(t.Date) == 0 || t.Date == "" {
		t.Date = todayDateStr
		return true, "", newTask
	}

	date, err := time.Parse(dataFormat, t.Date)
	if err != nil {
		return false, err.Error(), newTask
	}

	d := (24 * time.Hour)

	if date.Before(time.Now().Truncate(d)) {
		ruleIsSet := !(len(t.Repeat) == 0 || t.Repeat == "")
		if !ruleIsSet {
			t.Date = todayDateStr
			return true, todayDateStr, newTask
		}
		if ruleIsSet {
			newDateStr, err := NextDate(time.Now(), t.Date, t.Repeat)

			if err != nil {
				return false, err.Error(), newTask
			}
			t.Date = newDateStr
			return true, "", newTask
		}
	}
	return true, "", newTask
}

func (t *Task) IsValid() (bool, string) {
	if len(t.Title) == 0 || t.Title == "" {
		return false, "no title"
	}
	return true, ""
}

func CreateDB(DBFile string) {
	db, err := sql.Open("sqlite3", DBFile)

	if err != nil {
		return
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"date" CHAR(8) NOT NULL,
			"title" VARCHAR(255) NOT NULL,
			"comment" TEXT,
			"repeat" VARCHAR(128) NOT NULL
			)
		`)
	if err != nil {
		return
	}
	_, err = db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_date ON sсheduler (date)
		`)
	if err != nil {
		return
	}
}

func OpenDB(DBFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DBFile)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	config := &Config{}

	config.port = portConf
	config.dBFile = dBFileConf

	if len(os.Getenv("TODO_PORT")) > 0 {
		config.port = os.Getenv("TODO_PORT")
	}

	appPath, err := os.Executable()
	if err != nil {
		log.Println(err)
		return
	}

	if len(os.Getenv("TODO_DBFILE")) > 0 {
		config.dBFile = filepath.Join(filepath.Dir(appPath), "TODO_DBFILE")
	}
	err = ServerStart(config.port, config.dBFile)
	if err != nil {
		log.Println(err)
		return
	}
}
