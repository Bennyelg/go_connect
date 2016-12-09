package main

import (
	"fmt"
)

func main() {

	connection := new(Database)
	connection.GetDatabaseDetails("production")
	cursor := connection.Connect()
	defer cursor.Close()
	rows, err := cursor.Query("SELECT CURRENT date")
	if err != nil{
		fmt.Println(err)
	}
	for rows.Next() {
		var time string
		rows.Scan(&time)
		fmt.Println(time)
	}
}

