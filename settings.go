package main

import "database/sql"

const webDir = "./web"
const dBFileConf = "./scheduler.db"

var db *sql.DB

const dataFormat = "20060102"
const tasksLimit = 10

const portConf = ":7540"
