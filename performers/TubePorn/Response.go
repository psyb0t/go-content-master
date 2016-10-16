package TubePorn

import (
    "encoding/json"
)

func (p Performer) ErrorResponse(message string) error {
    return json.NewEncoder(p.Ctx).Encode(&Response{
        Status: "ERROR",
        Message: message,
    })
}

func (p Performer) OkResponse(message string) error {
    response := &Response{}
    response.Status = "OK"
    response.Message = message

    return json.NewEncoder(p.Ctx).Encode(response)
}

func (p Performer) OkVideoResponse(message string, data *Video) error {
    response := &VideoResponse{}
    response.Status = "OK"
    response.Message = message
    response.Data = data

    return json.NewEncoder(p.Ctx).Encode(response)
}

func (p Performer) OkVideosResponse(message string, data *Videos) error {
    response := &VideosResponse{}
    response.Status = "OK"
    response.Message = message
    response.Data = data

    return json.NewEncoder(p.Ctx).Encode(response)
}

func (p Performer) OkCategoryResponse(message string, data *Category) error {
    response := &CategoryResponse{}
    response.Status = "OK"
    response.Message = message
    response.Data = data

    return json.NewEncoder(p.Ctx).Encode(response)
}

func (p Performer) OkCategoriesResponse(message string,
  data *Categories) error {
    response := &CategoriesResponse{}
    response.Status = "OK"
    response.Message = message
    response.Data = data

    return json.NewEncoder(p.Ctx).Encode(response)
}
