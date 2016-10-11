package TubePorn

import (
    "fmt"
    "strconv"
    "encoding/json"

    "gopkg.in/redis.v4"
)

func (p Performer) Do(params []string) error {
    if len(params) < 2 {
        return p.ErrorResponse("No method specified")
    }

    p.Redis = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })

    request_method := string(p.Request.Method)
    switch params[1] {
        case "fetchvideo":
            if request_method == "GET" {
                return p.GetVideo(params)
            }

            if request_method == "POST" {
                return p.AddVideo(params)
            }

        case "fetchvideos":
            return p.GetVideos(params)

        default:
            return p.ErrorResponse("Invalid method")
    }

    return nil
}

func (p Performer) RKey(key_part string) string {
    full_key := fmt.Sprintf("%s_%s", p.KeyPrefix, key_part)
    return full_key
}

func (p Performer) GetVideo(params []string) error {
    if len(params) < 3 {
        return p.ErrorResponse("No item specified")
    }

    seo_title := params[2]

    redis_key := p.RKey(fmt.Sprintf("video:%s", seo_title))
    redis_res, err := p.Redis.Get(redis_key).Result()

    if err == redis.Nil {
        return p.ErrorResponse("Video does not exist")
    }

    if err != nil {
        return p.ErrorResponse("DB statement error")
    }

    video := &Video{SeoTitle: seo_title}
    err = json.Unmarshal([]byte(redis_res), &video)

    if err != nil {
        return p.ErrorResponse("Error decoding video")
    }

    return p.OkVideoResponse("Video fetched", video)
}

func (p Performer) AddVideo(data []string) error {
    return p.OkResponse("Added Video")
}

func (p Performer) GetVideos(data []string) error {
    if len(data) < 3 {
        return p.ErrorResponse("Invalid params")
    }

    number, err := strconv.Atoi(data[2])

    if err != nil {
        return p.ErrorResponse("Invalid number param")
    }

    offset := 0
    if len(data) > 3 {
        offset, err = strconv.Atoi(data[3])

        if err != nil {
            return p.ErrorResponse("Invalid number param")
        }
    }

    //category := ""
    //if len(data) > 4 {
    //    category = data[4]
    //}

    start_pos := offset
    end_pos := start_pos + number - 1

    videos_raw, err := p.Redis.ZRevRange(p.RKey("videos"),
                    int64(start_pos), int64(end_pos)).Result()

    if err == redis.Nil {
        return p.ErrorResponse("No videos found")
    }

    if err != nil {
        return p.ErrorResponse("DB statement error")
    }

    videos := &Videos{}
    for _, video_raw := range videos_raw {
        video := &Video{}
        err = json.Unmarshal([]byte(video_raw), &video)

        if err != nil {
            continue
        }

        *videos = append(*videos, video)
    }

    return p.OkVideosResponse("Videos fetched", videos)
}
