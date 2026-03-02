package main

import (
    "fmt"       // Package for printing text to console or writing responses
    "net/http"  // Package for creating HTTP servers and handling requests
)

func main() {
    capacity := 100
    port := 8080

    c := newCache(capacity)
    s := newServer(c)

    http.handleFunc("/Get", s.GetHandler)
    http.handleFunc("/Set", s.SetHandler)
    http.handleFunc("/Delete", s.DeleteHandler)

    http.ListenAndSever(port)
}