package main

import (
    "log"
    "encoding/json"

    "github.com/psyb0t/go-sfo"
)

var config_file = "/etc/content-master/config.json"

type Config struct {
    ListenHost string
    ListenPort int
}

var config *Config

func SetupConfig() {
    ConfigFromFile()

    if config != nil {
        return
    }

    config = &Config{
        ListenHost: "127.0.0.1",
        ListenPort: 60651,
    }

    SaveConfig()
}

func ConfigFromFile() bool {
    file, err := sfo.ReadFile(&config_file)

    if err != nil {
        log.Fatal(err)
    }

    err = json.Unmarshal(file.Bytes, config)

    if err != nil {
        return false
    }

    return true
}

func SaveConfig() bool {
    bytes, err := json.MarshalIndent(&config, "", "  ")

    if err != nil {
        log.Fatal(err)
    }

    err = sfo.WriteBytesToFile(&config_file, &bytes)

    if err != nil {
        return false
    }

    return true
}
