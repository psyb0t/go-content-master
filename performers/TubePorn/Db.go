package TubePorn

import (
    "fmt"
    "encoding/json"

    "gopkg.in/redis.v4"
)

func (p Performer) DbGetVideo(seo_title string, get_related bool) (*Video, error) {
    redis_key := p.RKey(fmt.Sprintf("video:%s", seo_title))
    redis_res, err := p.Redis.Get(redis_key).Result()

    if err != nil {
        if err != redis.Nil {
            return nil, err
        }
    }

    video := &Video{}
    err = json.Unmarshal([]byte(redis_res), &video)

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
    redis_res, err := p.Redis.ZRevRange(p.RKey("videos"),
        int64(start_pos), int64(end_pos)).Result()

    if err != nil {
        if err != redis.Nil {
            return nil, err
        }
    }

    videos := &Videos{}
    for _, video_raw := range redis_res {
        video := &Video{}
        err = json.Unmarshal([]byte(video_raw), &video)

        if err != nil {
            continue
        }

        *videos = append(*videos, video)
    }

    return videos, nil
}

func (p Performer) DbGetCategory(seo_title string) (*Category, error) {
    redis_res, err := p.Redis.Get(p.RKey(
        fmt.Sprintf("category:%s", seo_title))).Result()

    if err != nil {
        if err != redis.Nil {
            return nil, err
        }
    }

    category := &Category{}

    err = json.Unmarshal([]byte(redis_res), &category)

    if err != nil {
        return nil, err
    }

    return category, nil
}

func (p Performer) DbGetCategoryVideos(seo_title string,
  start_pos int, end_pos int) (*Videos, error) {
    redis_res, err := p.Redis.ZRevRange(p.RKey(
        fmt.Sprintf("category:%s:videos", seo_title)),
        int64(start_pos), int64(end_pos)).Result()

    if err != nil {
        if err != redis.Nil {
            return nil, err
        }
    }

    videos := &Videos{}
    for _, video_raw := range redis_res {
        video := &Video{}
        err = json.Unmarshal([]byte(video_raw), &video)

        if err != nil {
            continue
        }

        *videos = append(*videos, video)
    }

    return videos, nil
}

func (p Performer) DbGetCategories() (*Categories, error) {
    redis_resp, err := p.Redis.ZRevRange(p.RKey("categories"), 0, -1).Result()

    if err != nil {
        if err != redis.Nil {
            return nil, err
        }
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

    return categories, nil
}

func (p Performer) DbGetVideoSearch(kword string, start_pos int,
  end_pos int) (*Videos, error) {
    redis_res, err := p.Redis.Keys(p.RKey(
        fmt.Sprintf("video:*%s*", kword))).Result()

    if err != nil {
        if err != redis.Nil {
            return nil, err
        }
    }

    videos := &Videos{}
    for _, video_key := range redis_res {
        video := &Video{}

        video_raw, err := p.Redis.Get(video_key).Result()

        if err != nil {
            if err != redis.Nil {
                return nil, err
            }
        }

        err = json.Unmarshal([]byte(video_raw), &video)

        if err != nil {
            continue
        }

        *videos = append(*videos, video)
    }

    return videos.Range(start_pos, end_pos), nil
}
