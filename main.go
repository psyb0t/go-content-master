package main

import (
    "fmt"
    "strings"
    "time"

    "github.com/valyala/fasthttp"
    "github.com/garyburd/redigo/redis"
    "go-content-master/performers/TubePorn"
)

var redis_pool = &redis.Pool{
    MaxIdle: 3,
    IdleTimeout: 240 * time.Second,
    Dial: func () (redis.Conn, error) {
        c, err := redis.Dial("tcp", "127.0.0.1:6379")

        if err != nil {
            return nil, err
        }

        return c, err
    },
    TestOnBorrow: func(c redis.Conn, t time.Time) error {
        if time.Since(t) < time.Minute {
            return nil
        }
        _, err := c.Do("PING")
        return err
    },
}

func init() {
    config = nil
    SetupConfig()
    Log(fmt.Sprintf("Server started (%s:%d)",
        config.ListenHost, config.ListenPort), false)
}

func perform(ctx *fasthttp.RequestCtx) {
    Log(fmt.Sprintf("%s [%s] %s", strings.Split(
        ctx.RemoteAddr().String(), ":")[0],
        string(ctx.Method()), string(ctx.Path())), false)

    ctx.SetContentType("application/json")

    path := strings.Trim(string(ctx.Path()), "/")
    params := strings.Split(strings.ToLower(path), "/")

    if params[0] == "" {
        return
    }

    switch params[0] {
        case "tubeporn":
            performer := &TubePorn.Performer{
                KeyPrefix: "tubeporn",
                DBPool: redis_pool,
                Ctx: ctx,
            }

            performer.Do(params)
            break
    }
}

func main() {
    fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", config.ListenHost,
        config.ListenPort), perform)
}
