package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/microsoft/go-mssqldb"
)

type Database struct {
	ID   int64
	Name string
}

type Config struct {
	Server   string
	User     string
	Password string
	Port     int64
	LogFile  string
}

func main() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		panic(err.Error())
	}
	var config Config
	json.Unmarshal(file, &config)
	logger, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	defer logger.Close()
	log.SetOutput(logger)
	connectStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", config.Server, config.User, config.Password, config.Port)
	con, err := sql.Open("mssql", connectStr)
	if err != nil {
		log.Println("Connection failed ", err.Error())
	}
	defer con.Close()

	stmt := "SELECT database_id, name FROM sys.databases;"
	rows, err := con.Query(stmt)
	if err != nil {
		log.Println("Query failed ", err.Error())
	}
	for rows.Next() {
		var db Database
		err := rows.Scan(&db.ID, &db.Name)
		if err != nil {
			log.Println("Scan failed", err.Error())
		}
		_, err = con.Exec("INSERT INTO SOFJOBS.dbo.SampleDatabases (ID, DbName) VALUES (?, ?);", db.ID, db.Name)
		if err != nil {
			log.Println("Insert failed ", err.Error())
		}
	}
}
