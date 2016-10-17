package TubePorn

import (
    "fmt"
    "errors"
    "encoding/json"

    "github.com/garyburd/redigo/redis"
)

func (p Performer) DbGetVideo(seo_title string, get_related bool) (*Video, error) {
    redis_key := p.RKey(fmt.Sprintf("video:%s", seo_title))
    redis_res, err := redis.Bytes(p.Redis.Do("GET", redis_key))

    if err != nil {
        return nil, err
    }

    video := &Video{}
    err = json.Unmarshal(redis_res, &video)

    if err != nil {
        return nil, err
    }

    if get_related {
        for _, category := range video.Categories {
            category_videos, err := p.DbGetCategoryVideos(
                category.SeoTitle, 0, 10)

            if err != nil {
                continue
            }

            category.Videos = *category_videos
        }
    }

    return video, nil
}

func (p Performer) DbGetVideos(start_pos int, end_pos int) (*Videos, error) {
    redis_res, err := redis.ByteSlices(p.Redis.Do("ZREVRANGE",
        p.RKey("videos"), start_pos, end_pos))

    if err != nil {
        return nil, err
    }

    videos := &Videos{}
    for _, video_raw := range redis_res {
        video := &Video{}
        err = json.Unmarshal(video_raw, &video)

        if err != nil {
            continue
        }

        *videos = append(*videos, video)
    }

    return videos, nil
}

func (p Performer) DbGetCategory(seo_title string) (*Category, error) {
    redis_res, err := redis.Bytes(p.Redis.Do("GET", p.RKey(
        fmt.Sprintf("category:%s", seo_title))))

    if err != nil {
        return nil, err
    }

    category := &Category{}

    err = json.Unmarshal(redis_res, &category)

    if err != nil {
        return nil, err
    }

    return category, nil
}

func (p Performer) DbGetCategoryVideos(seo_title string,
  start_pos int, end_pos int) (*Videos, error) {
    redis_res, err := redis.ByteSlices(p.Redis.Do("ZREVRANGE",
        p.RKey(fmt.Sprintf("category:%s:videos", seo_title)),
        start_pos, end_pos))

    if err != nil {
        return nil, err
    }

    videos := &Videos{}
    for _, video_raw := range redis_res {
        video := &Video{}
        err = json.Unmarshal(video_raw, &video)

        if err != nil {
            continue
        }

        *videos = append(*videos, video)
    }

    return videos, nil
}

func (p Performer) DbGetCategories() (*Categories, error) {
    redis_resp, err := redis.ByteSlices(p.Redis.Do(
        "ZREVRANGE", p.RKey("categories"), 0, -1))

    if err != nil {
        return nil, err
    }

    categories := &Categories{}
    for _, category_raw := range redis_resp {
        category := &Category{}
        err = json.Unmarshal(category_raw, &category)

        if err != nil {
            continue
        }

        *categories = append(*categories, category)
    }

    return categories, nil
}

func (p Performer) DbGetVideoSearch(kword string, start_pos int,
  end_pos int) (*Videos, error) {
    redis_res, err := redis.ByteSlices(p.Redis.Do(
        "KEYS", p.RKey(fmt.Sprintf("video:*%s*", kword))))

    if err != nil {
        return nil, err
    }

    videos := &Videos{}
    for _, video_key := range redis_res {
        video := &Video{}

        video_raw, err := redis.Bytes(p.Redis.Do("GET", string(video_key)))

        if err != nil {
            continue
        }

        err = json.Unmarshal(video_raw, &video)

        if err != nil {
            continue
        }

        *videos = append(*videos, video)
    }

    return videos.Range(start_pos, end_pos), nil
}

func (p Performer) DbAddVideo(video *Video) error {
    video_json, err := json.Marshal(video)

    if err != nil {
        return err
    }

    _, err = p.DbGetVideo(video.SeoTitle, false)

    if err == nil {
        return errors.New("Video already exists")
    }

    p.DbSize = p.DbSize + 1

    _, err = p.Redis.Do("ZADD", p.RKey("videos"),
        float64(p.DbSize), video_json)

    p.DbSize = p.DbSize + 1

    video_key := p.RKey(fmt.Sprintf("video:%s", video.SeoTitle))
    _, err = p.Redis.Do("SET", video_key, video_json)

    if err != nil {
        return err
    }

    for i:=0; i<video.Categories.Length(); i++ {
        err := p.DbAddCategory(video.Categories[i])

        if err != nil {
            return err
        }

        cat_vids_key := p.RKey(fmt.Sprintf(
            "category:%s:videos", video.Categories[i].SeoTitle))

        p.DbSize = p.DbSize + 1

        _, err = p.Redis.Do("ZADD", cat_vids_key,
            float64(p.DbSize), video_json)

        if err != nil {
            return err
        }

    }

    return nil
}

func (p Performer) DbAddCategory(category *Category) error {
    category_json, err := json.Marshal(category)

    p.DbSize = p.DbSize + 1

    _, err = p.Redis.Do("ZADD", p.RKey("categories"), p.DbSize, category_json)

    p.DbSize = p.DbSize + 1

    category_key := p.RKey(fmt.Sprintf("category:%s", category.SeoTitle))
    _, err = p.Redis.Do("SET", category_key, category_json)

    if err != nil {
        return err
    }

    return nil
}
