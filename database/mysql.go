package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//log "github.com/sirupsen/logrus"
)
var (
	dbUser = "wunan"
	dbPass = "z0F5gb8WUhSHdTZ1"
	dbHost = "master.ccub4ta9krao.ap-southeast-2.rds.amazonaws.com"
	dbPort = "3306"
	dbName = "eiz_test"
	taskDbName = "eiz_task_test"
	db3plName = "eiz_3PL"
)

func InitializeDB() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	conn, err := sql.Open("mysql", dataSourceName)
	//defer conn.Close()
	if err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	//db = conn
	return conn, nil
}

func InitializeTaskDB() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, taskDbName)
	conn, err := sql.Open("mysql", dataSourceName)
	//defer conn.Close()
	if err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	//taskDb = conn
	return conn, nil
}

func Initialize3PlDB() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, db3plName)
	conn, err := sql.Open("mysql", dataSourceName)
	//defer conn.Close()
	if err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	return conn, nil
	//db3pl = conn
}
/*func GetDB() (*sql.DB, error) {
	once.Do(initializeDB)
	fmt.Println(123123)
	return db, nil
}*/

/*func GetTaskDB() (*sql.DB, error) {

	once.Do(initializeTaskDB)
	fmt.Println(taskDb)
	return taskDb, nil
}*/

/*func Get3PlDB() (*sql.DB, error) {
	once.Do(initialize3PlDB)
	fmt.Println(db3pl)
	return db3pl, nil
}*/