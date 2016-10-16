package TubePorn

import (
    "fmt"
    "strconv"
    "regexp"
    "strings"
    "encoding/json"

    "github.com/garyburd/redigo/redis"
)

func (p Performer) Do(params []string) error {
    var err error

    if len(params) < 2 {
        return p.ErrorResponse("No method specified")
    }

    p.Redis = p.DBPool.Get()
    defer p.Redis.Close()

    if err != nil {
        return p.ErrorResponse("Could not connect to the database")
    }

    p.DbSize, err = redis.Int(p.Redis.Do("DBSIZE"))

    if err != nil {
        return p.ErrorResponse("DBSIZE error")
    }

    request_method := string(p.Ctx.Method())
    switch params[1] {
        case "video":
            if request_method == "GET" {
                return p.GetVideo(params)
            }

            if request_method == "POST" {
                return p.AddVideo()
            }

        case "videos":
            return p.GetVideos(params)

        case "category":
            return p.GetCategory(params)

        case "categories":
            return p.GetCategories()

        case "random-video":
            return p.GetRandomVideo()

        case "search-video":
            return p.GetVideoSearch(params)

        default:
            return p.ErrorResponse("Invalid method")
    }

    p.Redis.Close()

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

    video, err := p.DbGetVideo(seo_title, true)

    if err != nil {
        return p.ErrorResponse("Could not get video")
    }

    return p.OkVideoResponse("Video fetched", video)
}

func (p Performer) AddVideo() error {
    video := &Video{}
    err := json.Unmarshal(p.Ctx.PostBody(), video)

    if err != nil {
        return p.ErrorResponse("Invalid video data")
    }

    err = p.DbAddVideo(video)

    if err != nil {
        return p.ErrorResponse("Could not add video")
    }

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

    start_pos := offset
    end_pos := start_pos + number - 1

    videos, error := p.DbGetVideos(start_pos, end_pos)

    if error != nil {
        return p.ErrorResponse("Could not get videos")
    }

    return p.OkVideosResponse("Videos fetched", videos)
}

func (p Performer) GetCategory(params []string) error {
    var err error

    if len(params) < 3 {
        return p.ErrorResponse("Invalid params")
    }

    seo_title := params[2]

    number := 50
    if len(params) > 3 {
        number, err = strconv.Atoi(params[3])

        if err != nil {
            return p.ErrorResponse("Invalid number param")
        }
    }

    offset := 0
    if len(params) > 4 {
        offset, err = strconv.Atoi(params[4])

        if err != nil {
            return p.ErrorResponse("Invalid number param")
        }
    }

    start_pos := offset
    end_pos := start_pos + number - 1

    category, err := p.DbGetCategory(seo_title)

    if err != nil {
        return p.ErrorResponse("Could not get category")
    }

    category_videos, _ := p.DbGetCategoryVideos(
        seo_title, start_pos, end_pos)

    category.Videos = *category_videos

    return p.OkCategoryResponse("Category fetched", category)
}

func (p Performer) GetCategories() error {
    categories, err := p.DbGetCategories()

    if err != nil {
        return p.ErrorResponse("Could not get categories")
    }

    return p.OkCategoriesResponse("Categories fetched", categories)
}

func (p Performer) GetVideoSearch(params []string) error {
    var err error

    if len(params) < 3 {
        return p.ErrorResponse("Invalid params")
    }

    kword := strings.Replace(params[2], " ", "-", -1)

    r, _ := regexp.Compile("[^0-9a-z-]")
    kword = r.ReplaceAllString(kword, "")

    r, _ = regexp.Compile("([-]+)")
    kword = r.ReplaceAllString(kword, "-")

    if len(kword) < 2 {
        return p.ErrorResponse("Invalid keyword length")
    }

    number := 50
    if len(params) > 3 {
        number, err = strconv.Atoi(params[3])

        if err != nil {
            return p.ErrorResponse("Invalid number param")
        }
    }

    offset := 0
    if len(params) > 4 {
        offset, err = strconv.Atoi(params[4])

        if err != nil {
            return p.ErrorResponse("Invalid number param")
        }
    }

    start_pos := offset
    end_pos := start_pos + number

    videos, err := p.DbGetVideoSearch(kword, start_pos, end_pos)

    if err != nil {
        return p.ErrorResponse("Could not get search results")
    }

    return p.OkVideosResponse("Got video search results", videos)
}

func (p Performer) GetRandomVideo() error {
    videos, err := p.DbGetVideos(0, 20)

    if err != nil {
        return p.ErrorResponse("Could not get a random video")
    }

    return p.OkVideoResponse("Random video fetched", videos.Rand())
}
