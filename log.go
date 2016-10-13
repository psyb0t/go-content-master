package main

import (
    "os"
    "log"

    "github.com/psyb0t/go-sfo"
)

var log_file = "/var/log/content-master/content-master.log"

func OpenLogFile() *os.File {
    err := sfo.CreateFile(&log_file)

    if err != nil {
        log.Fatal(err)
    }

    f, err := os.OpenFile(log_file,
        os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)


    if err != nil {
        log.Fatal(err)
    }

    return f
}
