package main

import (
    "strings"
    "net/http"
    "content-master/TubePorn"
)

func perform(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    path := strings.Trim(r.URL.Path, "/")
    params := strings.Split(strings.ToLower(path), "/")

    if params[0] == "" {
        return
    }

    switch params[0] {
        case "tubeporn":
            performer := &TubePorn.Performer{}
            performer.KeyPrefix = "tubeporn"
            performer.RespWriter = w
            performer.Request = r
            performer.Do(params)

            break
    }
}

func main() {
    http.HandleFunc("/", perform)
    http.ListenAndServe("127.0.0.1:8080", nil)
}
