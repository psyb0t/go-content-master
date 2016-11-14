package main

import (
    "sync"
    "time"
    "strings"
    "encoding/json"
    "fmt"
    "math/rand"

    "github.com/garyburd/redigo/redis"
    "github.com/psyb0t/go-fdb"
)

var err error

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

type Category struct {
    SeoTitle string `json:"seo_title"`
    Title string `json:"title"`
    Videos `json:"videos"`
}

type Categories []*Category
func (c Categories) Length() int {
    return len(c)
}
func (c Categories) Rand() *Category {
    return c[rand.Intn(len(c))]
}

type Video struct {
    SiteId string `json:"site_id"`
    Title string `json:"title"`
    SeoTitle string `json:"seo_title"`
    Description string `json:"description"`
    Thumbnail string `json:"thumbnail"`
    EmbedCode string `json:"embed_code"`
    Timestamp int `json:"timestamp"`
    Categories `json:"categories"`
}

type Videos []*Video
func (v Videos) Length() int {
    return len(v)
}
func (v Videos) Rand() *Video {
    return v[rand.Intn(len(v))]
}
func (v Videos) Range(start int, end int) *Videos {
    if v.Length() < start {
        start = v.Length()
    }

    if v.Length() < end {
        end = v.Length()
    }

    range_videos := v[start:end]
    return &range_videos
}
func (v Videos) AppendVideo(video *Video) *Videos {
    v = append(v, video)

    return &v
}

func main() {
    var wg sync.WaitGroup

    wg.Add(1)
    go func (wg *sync.WaitGroup) {
        rediss := redis_pool.Get()
        fdb_collection, _ := fdb.NewCollection(
            "./db/TubePorn/Videos")

        redis_res, _ := redis.ByteSlices(rediss.Do("KEYS", "tubeporn_video:*"))

        for _, keyname := range redis_res {
            value, _ := redis.Bytes(rediss.Do("GET", keyname))

            keyname = []byte(strings.Replace(string(keyname), "tubeporn_video:", "", -1))

            fdb_collection.Set(string(keyname), string(value))
        }

        fmt.Println(len(redis_res), "Videos Done")
        wg.Done()
    }(&wg)

    wg.Add(1)
    go func (wg *sync.WaitGroup) {
        rediss := redis_pool.Get()
        fdb_collection, _ := fdb.NewCollection(
            "./db/TubePorn/Categories")

        redis_res, _ := redis.ByteSlices(rediss.Do("KEYS", "tubeporn_category:*"))

        for _, keyname := range redis_res {
            value, _ := redis.Bytes(rediss.Do("GET", keyname))

            if strings.Contains(string(keyname), ":videos") {
                continue
            }

            keyname = []byte(strings.Replace(string(keyname), "tubeporn_category:", "", -1))

            fdb_collection.Set(string(keyname), string(value))
        }

        fmt.Println(len(redis_res), "Categories Done")
        wg.Done()
    }(&wg)


    wg.Add(1)
    go func (wg *sync.WaitGroup) {
        rediss := redis_pool.Get()

        redis_res, _ := redis.ByteSlices(rediss.Do("KEYS", "tubeporn_category:*:videos"))

        fdb_collection, _ := fdb.NewCollection(
            "./db/TubePorn/CategoryVideos")

        for _, keyname := range redis_res {
            catname := strings.Replace(
                strings.Replace(string(keyname), "tubeporn_category:", "", -1),
                ":videos", "", -1)

            vid_raws, _ := redis.ByteSlices(rediss.Do("ZRANGE", keyname, 0, -1))

            vid_list := []string{}
            for _, video_raw := range vid_raws {
                video := &Video{}

                err = json.Unmarshal(video_raw, &video)

                if err != nil {
                    continue
                }

                vid_list = append(vid_list, video.SeoTitle)
            }

            vid_list_json, err := json.Marshal(vid_list)

            if err != nil {
                continue
            }

            fdb_collection.Set(catname, string(vid_list_json))
        }

        wg.Done()
    }(&wg)



    wg.Wait()
}
