package main

import (
    "fmt"
    "strings"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "go-content-master/modules/TubePorn"
)

type Config struct {
    ListenHost string
    ListenPort int
}

var config *Config

func init() {
    config = &Config{}
    config.ListenHost = "127.0.0.1"
    config.ListenPort = 60651

    file, err := ioutil.ReadFile("./config.json")

    if err != nil {
        return
    }

    err = json.Unmarshal(file, &config)

    if err != nil {
        log.Fatal(err)
    }
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
