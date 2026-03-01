package main

import (
    "fmt"       // Package for printing text to console or writing responses
    "net/http"  // Package for creating HTTP servers and handling requests
)

func main() {
    c := cache.New(capacity)

    server := server.new(c)
    server.Start(":8080")
}