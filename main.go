package main

import (
    "fmt"
    "strings"
    "net/http"
    "runtime/debug"

    "go-content-master/performers/TubePorn"
)

func init() {
    config = nil
    SetupConfig()
    Log(fmt.Sprintf("Server started (%s:%d)",
        config.ListenHost, config.ListenPort), false)
}

func perform(w http.ResponseWriter, r *http.Request) {
    defer debug.FreeOSMemory()

    Log(fmt.Sprintf("%s [%s] %s", strings.Split(r.RemoteAddr, ":")[0],
        string(r.Method), r.URL.Path), false)

    w.Header().Set("Content-Type", "application/json")

    path := strings.Trim(r.URL.Path, "/")
    params := strings.Split(strings.ToLower(path), "/")

    if params[0] == "" {
        return
    }

    switch params[0] {
        case "tubeporn":
            performer := &TubePorn.Performer{
                KeyPrefix: "tubeporn",
                RespWriter: w,
                Request: r,
            }

            performer.Do(params)
            break
    }
}

func main() {
    http.HandleFunc("/", perform)
    http.ListenAndServe(fmt.Sprintf("%s:%d", config.ListenHost,
        config.ListenPort), nil)
}
