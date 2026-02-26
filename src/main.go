package main

import (
    "fmt"       // Package for printing text to console or writing responses
    "net/http"  // Package for creating HTTP servers and handling requests
)

// helloHandler is a function that handles HTTP requests to the "/" path
// w is used to send a response back to the client (browser)
// r contains information about the client's request
func helloHandler(w http.ResponseWriter, r *http.Request) {
    // Write "Hello, World!" to the HTTP response body
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    // Register the handler function for the root path "/"
    // Whenever someone visits "/", Go will call helloHandler
    http.HandleFunc("/", helloHandler)

    // Print to the console that the server is starting
    fmt.Println("Server starting on port 8080...")

    // Start the server on port 8080
    // nil means we’re using the default ServeMux (default router)
    // ListenAndServe blocks and listens for incoming requests
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        // If there’s an error starting the server, print it
        fmt.Println("Error starting server:", err)
    }
}