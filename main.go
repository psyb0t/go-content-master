package main

import (
    "log"
    "fmt"
    "strings"

    "github.com/valyala/fasthttp"
    "github.com/psyb0t/go-fdb"

    "go-content-master/performers/TubePorn"
    "go-content-master/performers/MovieWatch"
)

var err error

var (
    TubePornFDB map[string]*fdb.Collection
    MovieWatchFDB map[string]*fdb.Collection
)

func init() {
    config = nil
    SetupConfig()
    Log(fmt.Sprintf("Server started (%s:%d)",
        config.ListenHost, config.ListenPort), false)

    TubePornFDB = make(map[string]*fdb.Collection)

    TubePornFDB["videos"], err = fdb.NewCollection(
        "/etc/fdb/TubePorn/Videos")

    if err != nil {
        log.Fatal("Could not create TubePorn:Videos collection")
    }

    TubePornFDB["categories"], err = fdb.NewCollection(
        "/etc/fdb/TubePorn/Categories")

    if err != nil {
        log.Fatal("Could not create TubePorn:Categories collection")
    }

    TubePornFDB["category_videos"], err = fdb.NewCollection(
        "/etc/fdb/TubePorn/CategoryVideos")

    if err != nil {
        log.Fatal("Could not create TubePorn:Categories collection")
    }

    MovieWatchFDB = make(map[string]*fdb.Collection)

    MovieWatchFDB["videos"], err = fdb.NewCollection(
        "/etc/fdb/MovieWatch/Videos")

    if err != nil {
        log.Fatal("Could not create MovieWatch:Videos collection")
    }

    MovieWatchFDB["genres"], err = fdb.NewCollection(
        "/etc/fdb/MovieWatch/Genres")

    if err != nil {
        log.Fatal("Could not create MovieWatch:Genres collection")
    }
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
                Ctx: ctx,
                FDBVideos: TubePornFDB["videos"],
                FDBCategories: TubePornFDB["categories"],
                FDBCategoryVideos: TubePornFDB["category_videos"],
            }

            performer.Do(params)
            break

        case "moviewatch":
            performer := &MovieWatch.Performer{
                Ctx: ctx,
                FDBVideos: MovieWatchFDB["videos"],
                FDBGenres: MovieWatchFDB["genres"],
            }

            performer.Do(params)
            break
    }
}



func main() {
    fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", config.ListenHost,
        config.ListenPort), perform)
}
