package main

import (
    "fmt"       // Package for printing text to console or writing responses
    "net/http"  // Package for creating HTTP servers and handling requests
)

func main() {
    capacity := 100
    port := 8080

    c := newCache(capacity)
    s := NewServer(c)

    http.HandleFunc("/Get", s.GetHandler)
    http.HandleFunc("/Set", s.SetHandler)
    http.HandleFunc("/Delete", s.DeleteHandler)

    http.ListenAndServe(port)

}