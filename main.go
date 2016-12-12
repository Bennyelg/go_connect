package main

import (
	"fmt"
	cnt "dbConnector"
)

const (
	databaseFilePath = "/Users/benny/go/database.toml"
	env = "production"
)

func main() {

	connection := new(cnt.Database)
	connection.ParseDatabaseByEnv(env, databaseFilePath)
	cursor := connection.Connect()
	defer cursor.Close()
	rows, err := cursor.Query("SELECT CURRENT_DATE")
	if err != nil{
		fmt.Println(err)
	}
	for rows.Next() {
		var time string
		rows.Scan(&time)
		fmt.Println(time)
	}
}
