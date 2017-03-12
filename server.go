// To run:
// $ go run server.go

package main

import (
  "fmt"
  "io/ioutil"
  "strconv"
  "log"
  "net/http"
  "encoding/json"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type Item struct {
  Id int
  Name string
}

// Handle all requests
func Handler(response http.ResponseWriter, request *http.Request) {
  response.Header().Set("Content-type", "text/html")
  webpage, err := ioutil.ReadFile("index.html")
  if err != nil {
    http.Error(response, fmt.Sprintf("home.html file error %v", err), 500)
  }
  fmt.Fprint(response, string(webpage));
}

// Respond to URLs of form /generic/
func APIHandler(response http.ResponseWriter, request *http.Request) {
  // Connect to database
  db, e := sql.Open("mysql", "username:password@tcp(localhost:3306)/farm")
    if(e != nil) {
      fmt.Print(e)
    }

    // set mime type to JSON
    response.Header().Set("Content-type", "application/json")
      err := request.ParseForm()
      if err != nil {
        http.Error(response, fmt.Sprintf("Error parsing url %v", err), 500)
      }

    // Unable to define dynamic slice in golang
    var result = make([]string, 1000)

    switch request.Method {
      case "GET":
        st, err := db.Prepare("select * from items limit 10")
        if err != nil {
          fmt.Print(err);
        }
        rows, err := st.Query()
        if err != nil {
          fmt.Print(err)
        }
        i := 0
        for rows.Next() {
          var name string
          var id int
          err = rows.Scan(&id, &name)
          item := &Item{Id: id, Name: name}
            b, err := json.Marshal(item)
            if err != nil {
              fmt.Println(err)
              return
            }
          result[i] = fmt.Sprintf("%s", string(b))
          i++
        }
        result = result[:i]

        case "POST":
          name := request.PostFormValue("name")
          st, err := db.Prepare("INSERT INTO items(name) VALUES(?)")
          if err != nil {
            fmt.Print(err);
          }
          res, err := st.Exec(name)
          if err != nil {
            fmt.Print(err)
          }
          if res != nil {
            result[0] = "true"
          }
          result = result[:1]

          default:
    }

    json, err := json.Marshal(result)
    if err != nil {
      fmt.Println(err)
      return
    }

    // Send text diagnostics to client
    fmt.Fprintf(response, "%v", string(json))
    db.Close()
}

func main() {
  port := 3000
  var err string
  portString := strconv.Itoa(port)

  mux := http.NewServeMux()
  mux.Handle("/api/", http.HandlerFunc(APIHandler))
  mux.Handle("/", http.HandlerFunc(Handler))

  // listen on given port with these routes on the server
  log.Print("Listening on port " + portString + "...")
  errs := http.ListenAndServe(":" + portString, mux)
  if errs != nil {
    log.Fatal("ListenAndServe erro: ", err)
  }
}
