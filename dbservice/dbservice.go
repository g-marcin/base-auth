package dbservice

import (
	"fmt"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB
var goquDb *goqu.Database

func init () {
	var err error
	db, err = sql.Open("postgres", "user=user password=pass dbname=user sslmode=disable")
	if err != nil {
		 panic(fmt.Sprintf("failed to ping database: %v", err))
	}
	goquDb = goqu.New("postgres", db)
}

func CheckCredentials(username, password string) (bool, error) {
  log.Println("Checking credentials for username:", username)

  ds := goquDb.From("users").Select(goqu.COUNT("*")).Where(goqu.Ex{
      "username": username,
      "password": password,
  })

  var count int
  found, err := ds.ScanVal(&count)
  if err != nil {
      log.Println("Error querying database:", err)
      return false, fmt.Errorf("query error: %v", err)
  }
  
  log.Println("Query result found:", found)
  log.Println("Matching rows count:", count)
  log.Println("Credentials valid:", count > 0)

  return count > 0, nil
}

func InsertUser (username, password string) error {
	ds := goquDb.Insert("users").Rows(
		goqu.Record{"username": username, "password": password},
	)
	
	sql, _, err := ds.ToSQL()
  if err != nil {
      return fmt.Errorf("failed to generate SQL: %v", err)
  }
	
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	} 
	
	return nil
}
