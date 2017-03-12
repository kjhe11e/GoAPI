// To run:
// $ go run server.go

package main

import (
  "fmt"
  "bytes"
  "database/sql"
  "net/http"

  "github.com/gin-gonic/gin"
  _ "github.com/go-sql-driver/mysql"
)

type Item struct {
  id int
  name string
}

func main() {
  db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/go_api_test")
  if err != nil {
    fmt.Print(err.Error())
  }
  defer db.Close()

  // make sure connection is available
  err = db.Ping()
  if err != nil {
    fmt.Print(err.Error())
  }

  router := gin.Default()

  // TODO: move API handlers out of main func
  // GET an Item
  router.GET("/item/:id", func(c *gin.Context) {
    var (
      item Item
      result gin.H
    )
    id := c.Param("id")
    row := db.QueryRow("SELECT id, name from item where id = ?;", id)
    err = row.Scan(&item.id, &item.name)
    if err != nil {
      // if no result send null
      result = gin.H {
        "result": nil,
        "count": 0,
      }
    } else {
      result = gin.H {
        "id": item.id,
        "name": item.name,
      }
    }
    c.JSON(http.StatusOK, result)
  })

  // GET all items
  router.GET("/items", func(c *gin.Context) {
    var (
      item Item
      items []Item
    )
    rows, err := db.Query("SELECT id, name from item;")
    if err != nil {
      fmt.Print(err.Error())
    }
    for rows.Next() {
      err = rows.Scan(&item.id, &item.name)
      items = append(items, item)
      if err != nil {
        fmt.Print(err.Error())
      }
    }
    defer rows.Close()
    c.JSON(http.StatusOK, gin.H {
      "result": items,
      "count": len(items),
    })
  })

  // POST new item
  router.POST("/item", func(c *gin.Context) {
    var buffer bytes.Buffer
    name := c.PostForm("name")
    stmt, err := db.Prepare("INSERT into item (name) values(?);")
    if err != nil {
      fmt.Print(err.Error())
    }
    _, err = stmt.Exec(name)

    if err != nil {
      fmt.Print(err.Error())
    }

    // Fastest way to append strings
    buffer.WriteString(name)
    defer stmt.Close()
    tmpName := buffer.String()
    c.JSON(http.StatusOK, gin.H {
      "message": fmt.Sprintf(" %s successfully created", tmpName),
    })
  })

  router.Run(":3000")
}
