package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "strings"
    "net/http"

    "go-content-master/performers/TubePorn"
)

func init() {
    config = nil
    SetupConfig()

    log.SetOutput(io.MultiWriter(os.Stdout, OpenLogFile()))
}

func perform(w http.ResponseWriter, r *http.Request) {
    log.Printf("%s %s", strings.Split(r.RemoteAddr, ":")[0], r.URL.Path)

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
