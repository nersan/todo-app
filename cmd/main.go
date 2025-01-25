package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, ToDo App!")
    })

    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
