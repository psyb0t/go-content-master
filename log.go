package main

import (
    "os"
    "io"
    "log"

    "github.com/psyb0t/go-sfo"
)

var log_file = "/var/log/content-master/content-master.log"

func Log(item interface{}, fatal bool) {
    err := sfo.CreateFile(&log_file)

    if err != nil {
        log.Fatal(err)
    }

    f, err := os.OpenFile(log_file, os.O_RDWR | os.O_APPEND, 0644)
    defer f.Close()


    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(io.MultiWriter(os.Stdout, f))

    if fatal {
        log.Fatal(item)
    }

    log.Println(item)
}
