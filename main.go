package main

import (
    "log"
    "fmt"
    "strings"
    //"time"
    //"encoding/json"
    //"sync"

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


    //var wg sync.WaitGroup
    //
    //wg.Add(1)
    //go func (wg *sync.WaitGroup) {
    //    rediss := redis_pool.Get()
    //    fdb_collection, _ := fdb.NewCollection(
    //        "./db/TubePorn/Videos")
    //
    //    redis_res, _ := redis.ByteSlices(rediss.Do("KEYS", "tubeporn_video:*"))
    //
    //    for _, keyname := range redis_res {
    //        value, _ := redis.Bytes(rediss.Do("GET", keyname))
    //
    //        keyname = []byte(strings.Replace(string(keyname), "tubeporn_video:", "", -1))
    //
    //        fdb_collection.Set(string(keyname), string(value))
    //    }
    //
    //    fmt.Println(len(redis_res), "Videos Done")
    //    wg.Done()
    //}(&wg)
    //
    //wg.Add(1)
    //go func (wg *sync.WaitGroup) {
    //    rediss := redis_pool.Get()
    //    fdb_collection, _ := fdb.NewCollection(
    //        "./db/TubePorn/Categories")
    //
    //    redis_res, _ := redis.ByteSlices(rediss.Do("KEYS", "tubeporn_category:*"))
    //
    //    for _, keyname := range redis_res {
    //        value, _ := redis.Bytes(rediss.Do("GET", keyname))
    //
    //        if strings.Contains(string(keyname), ":videos") {
    //            continue
    //        }
    //
    //        keyname = []byte(strings.Replace(string(keyname), "tubeporn_category:", "", -1))
    //
    //        fdb_collection.Set(string(keyname), string(value))
    //    }
    //
    //    fmt.Println(len(redis_res), "Categories Done")
    //    wg.Done()
    //}(&wg)
    //
    //
    //wg.Add(1)
    //go func (wg *sync.WaitGroup) {
    //    rediss := redis_pool.Get()
    //
    //    redis_res, _ := redis.ByteSlices(rediss.Do("KEYS", "tubeporn_category:*:videos"))
    //
    //    fdb_collection, _ := fdb.NewCollection(
    //        "./db/TubePorn/CategoryVideos")
    //
    //    for _, keyname := range redis_res {
    //        catname := strings.Replace(
    //            strings.Replace(string(keyname), "tubeporn_category:", "", -1),
    //            ":videos", "", -1)
    //
    //        vid_raws, _ := redis.ByteSlices(rediss.Do("ZRANGE", keyname, 0, -1))
    //
    //        vid_list := []string{}
    //        for _, video_raw := range vid_raws {
    //            video := &TubePorn.Video{}
    //
    //            err = json.Unmarshal(video_raw, &video)
    //
    //            if err != nil {
    //                continue
    //            }
    //
    //            vid_list = append(vid_list, video.SeoTitle)
    //        }
    //
    //        vid_list_json, err := json.Marshal(vid_list)
    //
    //        if err != nil {
    //            continue
    //        }
    //
    //        fdb_collection.Set(catname, string(vid_list_json))
    //    }
    //
    //    wg.Done()
    //}(&wg)
    //
    //
    //
    //wg.Wait()
    //
    //log.Fatal("dafuq")


    TubePornFDB = make(map[string]*fdb.Collection)

    TubePornFDB["videos"], err = fdb.NewCollection(
        "./db/TubePorn/Videos")

    if err != nil {
        log.Fatal("Could not create TubePorn:Videos collection")
    }

    TubePornFDB["categories"], err = fdb.NewCollection(
        "./db/TubePorn/Categories")

    if err != nil {
        log.Fatal("Could not create TubePorn:Categories collection")
    }

    TubePornFDB["category_videos"], err = fdb.NewCollection(
        "./db/TubePorn/CategoryVideos")

    if err != nil {
        log.Fatal("Could not create TubePorn:Categories collection")
    }

    MovieWatchFDB = make(map[string]*fdb.Collection)

    MovieWatchFDB["videos"], err = fdb.NewCollection(
        "./db/MovieWatch/Videos")

    if err != nil {
        log.Fatal("Could not create MovieWatch:Videos collection")
    }

    MovieWatchFDB["genres"], err = fdb.NewCollection(
        "./db/MovieWatch/Genres")

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
