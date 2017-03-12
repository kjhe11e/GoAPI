package main

import (
  "database/sql"
  "fmt"
  _ "github.com/go-sql-driver/mysql"
)

func main() {
  db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/go_api_test")
  if err != nil {
    fmt.Println(err.Error())
  }

  defer db.Close()

  // make sure connection is available
  err = db.Ping()
  if err != nil {
    fmt.Println(err.Error())
  }

  stmt, err := db.Prepare("CREATE TABLE item (id int NOT NULL AUTO_INCREMENT, name varchar(100), PRIMARY KEY (id));")
  if err != nil {
    fmt.Println(err.Error())
  }

  _, err = stmt.Exec()
  if err != nil {
    fmt.Println(err.Error())
  } else {
    fmt.Println("Item table successfully migrated...")
  }
}
