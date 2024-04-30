package tools

import (
	"database/sql"
	"financeSys/database"
	"sync"
)

var (
	db *sql.DB
	taskDb *sql.DB
	Db3Pl *sql.DB
	once1 sync.Once
	once2 sync.Once
	once3 sync.Once
)
func init()  {
	once1.Do(GetDB)
	once2.Do(GetTaskDB)
	once3.Do(Get3PlDB)
}

func GetDB() {
	dbRes, _  := database.InitializeDB()
	db = dbRes
}

func GetTaskDB() {
	dbRes, _ := database.InitializeTaskDB()
	taskDb =  dbRes
}

func Get3PlDB() {
	dbRes, _ := database.Initialize3PlDB()
	Db3Pl =  dbRes
}

func CloseDb()  {
	db.Close()
	taskDb.Close()
	Db3Pl.Close()
}