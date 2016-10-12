package TubePorn

import (
    "net/http"
    "gopkg.in/redis.v4"
)

type Performer struct {
    KeyPrefix string
    Redis *redis.Client
    RespWriter http.ResponseWriter
    Request *http.Request
}

type Category struct {
    SeoTitle string `json:"seo_title"`
    Title string `json:"title"`
    Videos `json:"videos"`
}

type Categories []*Category

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

type Response struct {
    Status string `json:"status"`
    Message string `json:"message"`
}

type VideoResponse struct {
    Status string `json:"status"`
    Message string `json:"message"`
    Data *Video `json:"data"`
}

type VideosResponse struct {
    Status string `json:"status"`
    Message string `json:"message"`
    Data *Videos `json:"data"`
}

type CategoryResponse struct {
    Status string `json:"status"`
    Message string `json:"message"`
    Data *Category `json:"data"`
}

type CategoriesResponse struct {
    Status string `json:"status"`
    Message string `json:"message"`
    Data *Categories `json:"data"`
}