package TubePorn

import (
    "encoding/json"
)

func (p Performer) ErrorResponse(message string) error {
    return json.NewEncoder(p.RespWriter).Encode(&Response{
        Status: "ERROR",
        Message: message,
    })
}

func (p Performer) OkResponse(message string) error {
    return json.NewEncoder(p.RespWriter).Encode(&Response{
        Status: "OK",
        Message: message,
    })
}

func (p Performer) OkVideoResponse(message string, data *Video) error {
    return json.NewEncoder(p.RespWriter).Encode(&VideoResponse{
        Status: "OK",
        Message: message,
        Data: data,
    })
}

func (p Performer) OkVideosResponse(message string, data *Videos) error {
    return json.NewEncoder(p.RespWriter).Encode(&VideosResponse{
        Status: "OK",
        Message: message,
        Data: data,
    })
}

func (p Performer) OkCategoryResponse(message string, data *Category) error {
    return json.NewEncoder(p.RespWriter).Encode(&CategoryResponse{
        Status: "OK",
        Message: message,
        Data: data,
    })
}

func (p Performer) OkCategoriesResponse(message string,
  data *Categories) error {
    return json.NewEncoder(p.RespWriter).Encode(&CategoriesResponse{
        Status: "OK",
        Message: message,
        Data: data,
    })
}
