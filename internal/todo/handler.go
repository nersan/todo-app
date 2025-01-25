package todo

import (
    "fmt"
    "net/http"
)

func TodoHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "This is the ToDo handler!")
}
