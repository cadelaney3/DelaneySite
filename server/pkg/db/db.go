package db

import (
	_ "github.com/denisenkom/go-mssqldb" // import for MS Azure DB
    "database/sql"
    "context"
    "log"
    "fmt"
)

type azureDBConn struct {
	server string
	port int
	user string
	password string
	database string
}

type postgresConn struct {
	host string
	port int 
	user string
	password string 
	dbname string
}

// var azureDB *sql.DB

func initAzureDB(keys map[string]map[string]string) *sql.DB {
	conn := azureDBConn{
		server: "delaneysite.database.windows.net",
		port: 1433,
		user: keys["AZURE_DB"]["USER"],
		password: keys["AZURE_DB"]["PASSWORD"],
		database: "delaneysite",
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		conn.server, conn.user, conn.password, conn.port, conn.database)

	// var err error

	// Create connection pool
	azureDB, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
    err = azureDB.PingContext(ctx)
    if err != nil {
        log.Fatal(err.Error())
    }
	fmt.Printf("Connected!\n")

	return azureDB
}

func connDB(keys map[string]map[string]string) {

	postgres := postgresConn{
		host: "localhost",
		port: 5432,
		user: keys["POSTGRES"]["USER"],
		password: keys["POSTGRES"]["PASSWORD"],
		dbname: "delaneysite",
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
	postgres.host, postgres.port, postgres.user, postgres.password, postgres.dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
	  panic(err)
	} 
	fmt.Println("Successfully connected!")
}