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

        case "fetchcategories":
            return p.GetCategories()

        case "randomvideo":
            return p.RandomVideo()

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

func (p Performer) AddVideo(params []string) error {
    return p.OkResponse("Added Video")
}

func (p Performer) GetVideos(params []string) error {
    if len(params) < 3 {
        return p.ErrorResponse("Invalid params")
    }

    number, err := strconv.Atoi(params[2])

    if err != nil {
        return p.ErrorResponse("Invalid number param")
    }

    offset := 0
    if len(params) > 3 {
        offset, err = strconv.Atoi(params[3])

        if err != nil {
            return p.ErrorResponse("Invalid number param")
        }
    }

    category_seo_title := ""
    if len(params) > 4 {
        category_seo_title = params[4]
    }

    start_pos := offset
    end_pos := start_pos + number - 1

    if category_seo_title == "" {
        videos, error := p.DbGetVideos(start_pos, end_pos)

        if error != nil {
            return p.ErrorResponse("Could not get videos")
        }

        return p.OkVideosResponse("Videos fetched", videos)
    }

    return p.GetCategory(params)
}

func (p Performer) GetCategory(params []string) error {
    if len(params) < 3 {
        return p.ErrorResponse("Invalid params")
    }

    number, err := strconv.Atoi(params[2])

    if err != nil {
        return p.ErrorResponse("Invalid number param")
    }

    offset := 0
    if len(params) > 3 {
        offset, err = strconv.Atoi(params[3])

        if err != nil {
            return p.ErrorResponse("Invalid number param")
        }
    }

    category_seo_title := ""
    if len(params) > 4 {
        category_seo_title = params[4]
    }

    start_pos := offset
    end_pos := start_pos + number - 1

    redis_get_res, err := p.Redis.Get(p.RKey(
        fmt.Sprintf("category:%s", category_seo_title))).Result()

    if err == redis.Nil {
        return p.ErrorResponse("Category does not exist")
    }

    if err != nil {
        return p.ErrorResponse("DB statement error")
    }

    category := &Category{SeoTitle: category_seo_title}
    err = json.Unmarshal([]byte(redis_get_res), &category)

    if err != nil {
        return p.ErrorResponse("Error decoding category")
    }

    redis_zrevrange_res, err := p.Redis.ZRevRange(p.RKey(
                fmt.Sprintf("category:%s:videos", category_seo_title)),
                int64(start_pos), int64(end_pos)).Result()

    if err == redis.Nil {
        return p.ErrorResponse("No videos found")
    }

    if err != nil {
        return p.ErrorResponse("DB statement error")
    }

    for _, video_raw := range redis_zrevrange_res {
        video := &Video{}
        err = json.Unmarshal([]byte(video_raw), &video)

        if err != nil {
            continue
        }

        category.Videos = append(category.Videos, video)
    }

    return p.OkCategoryResponse("Category fetched", category)
}

func (p Performer) GetCategories() error {
    redis_resp, err := p.Redis.ZRevRange(p.RKey("categories"), 0, -1).Result()

    if err == redis.Nil {
        return p.ErrorResponse("No videos found")
    }

    if err != nil {
        return p.ErrorResponse("DB statement error")
    }

    categories := &Categories{}
    for _, category_raw := range redis_resp {
        category := &Category{}
        err = json.Unmarshal([]byte(category_raw), &category)

        if err != nil {
            continue
        }

        *categories = append(*categories, category)
    }

    return p.OkCategoriesResponse("Categories fetched", categories)
}

func (p Performer) RandomVideo() error {
    videos, err := p.DbGetVideos(0, 20)

    if err != nil {
        return p.ErrorResponse("Could not fetch a random video")
    }

    return p.OkVideoResponse("Random video fetched", videos.Rand())
}
