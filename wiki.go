package main
import (
       "fmt"
       "net/http"
)

func handler(writer http.ResponseWriter, request *http.Request) {
     fmt.Fprintf(writer, "Hi there, I love %s!", request.URL.Path[1:])
}

func main() {
     http.HandleFunc("/", handler)
     http.ListenAndServe(":5000", nil)
}
