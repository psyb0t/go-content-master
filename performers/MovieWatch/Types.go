package MovieWatch

import (
    "math/rand"

    "github.com/psyb0t/go-fdb"
    "github.com/valyala/fasthttp"
)

type Performer struct {
    Ctx *fasthttp.RequestCtx
    FDBGenres *fdb.Collection
    FDBVideos *fdb.Collection
    FDBCategories *fdb.Collection
    FDBCategoryVideos *fdb.Collection
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
    Id string `json:"id"`
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

type Response struct {
    Status string `json:"status"`
    Message string `json:"message"`
}

type VideoResponse struct {
    Response
    Data *Video `json:"data"`
}

type VideosResponse struct {
    Response
    Data *Videos `json:"data"`
}

type CategoryResponse struct {
    Response
    Data *Category `json:"data"`
}

type CategoriesResponse struct {
    Response
    Data *Categories `json:"data"`
}
