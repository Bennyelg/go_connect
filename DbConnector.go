package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"strconv"
)

type Database struct {
	dbType string
	username string
	password string
	port string
	database  string
	host string
	maxIdleConnection int
	maxOpenConnection int
}

// Parsing database information by giving environment. (development/production)
func (db *Database) GetDatabaseDetails(env string) {
	// Load file assuming is in the parent directory.
	databaseInfo, _ := toml.LoadFile("database.toml")
	// Parsing env.
	dbEnv := fmt.Sprintf("databases.%s", env)
	if dbEnv == "" {
		fmt.Println("[Error]: The env provided: ", env, " Not found.")
		os.Exit(1)
	}
	// Fetching the information according the env provided.
	infoTree := databaseInfo.Get(dbEnv).(*toml.TomlTree)
	username := infoTree.Get("user")
	password := infoTree.Get("password")
	port := infoTree.Get("port")
	database := infoTree.Get("database")
	host := infoTree.Get("host")
	dbType := infoTree.Get("type")
	maxOpenConnection := infoTree.Get("max_open_connection")
	maxIdleConnection := infoTree.Get("max_idle_connection")
	// Check if anything is missing.
	if username == nil || password == nil || port == nil || host == nil{
		fmt.Println("[Error]: Parameters missing! Please validate your " +
			".TOML file (Hint: Required parameters are: username, password, port, host)")
		os.Exit(1)
	}
	db.dbType = dbType.(string)
	db.username = username.(string)
	db.password = password.(string)
	db.database = database.(string)
	db.host = host.(string)
	db.port = port.(string)
	mic := maxIdleConnection.(string)
	moc := maxOpenConnection.(string)
	db.maxIdleConnection, _ = strconv.Atoi(mic)
	db.maxOpenConnection, _ = strconv.Atoi(moc)

}

// Connect into the database.
func (db *Database) Connect() *sql.DB{
	url := ""

	if db.dbType == "mysql"{
		url = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.username, db.password, db.host, db.port, db.database)

	}else if db.dbType == "postgres" {
		url= fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", db.username, db.password, db.database)
	}else{
		fmt.Println("[Error]: Not supported database.")
	}
	DbCon, err := sql.Open(db.dbType, url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	DbCon.Ping()
	DbCon.SetMaxOpenConns(db.maxOpenConnection)
	DbCon.SetMaxIdleConns(db.maxIdleConnection)
	fmt.Println("[Status]: Connected!\nMaxOpenConnection:", db.maxOpenConnection,
		    "\nMaxIdleConnection:", db.maxIdleConnection)
	return DbCon
}
