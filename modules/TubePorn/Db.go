package TubePorn

import (
    "encoding/json"
)

func (p Performer) DbGetVideos(start_pos int, end_pos int) (*Videos, error) {
    redis_res, err := p.Redis.ZRevRange(p.RKey("videos"),
                int64(start_pos), int64(end_pos)).Result()

    if err != nil {
        return nil, err
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
