package TubePorn

import (
    "errors"
    "encoding/json"
    "sync"
)

func (p Performer) DbGetVideo(seo_title string,
  get_related bool) (*Video, error) {
    value, err := p.FDBVideos.Get(seo_title)

    if err != nil {
        return nil, err
    }

    video := &Video{}

    err = json.Unmarshal([]byte(value), &video)

    if err != nil {
        return nil, err
    }

    if get_related {
        var wg sync.WaitGroup

        for _, category := range video.Categories {
            wg.Add(1)
            go func(wg *sync.WaitGroup, category *Category) {
                category_videos, err := p.DbGetCategoryVideos(
                    category.SeoTitle, 0, 10)

                if err != nil {
                    wg.Done()
                    return
                }

                category.Videos = *category_videos

                wg.Done()
            }(&wg, category)
        }

        wg.Wait()
    }

    return video, nil
}

func (p Performer) DbGetVideos(start_pos int, end_pos int) (*Videos, error) {
    db_vids, err := p.FDBVideos.GetReverseRange(start_pos, end_pos)

    if err != nil {
        return nil, err
    }

    if err != nil {
        return nil, err
    }

    videos := &Videos{}
    for _, video_raw := range db_vids {
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
    categ_data, err := p.FDBCategories.Get(seo_title)

    if err != nil {
        return nil, err
    }

    category := &Category{}

    err = json.Unmarshal([]byte(categ_data), &category)

    if err != nil {
        return nil, err
    }

    return category, nil
}

func (p Performer) DbGetCategoryVideos(seo_title string,
  start_pos int, end_pos int) (*Videos, error) {
    vid_res := []string{}

    fdb_res, err := p.FDBCategoryVideos.Get(seo_title)

    if err != nil {
        return nil, err
    }

    err = json.Unmarshal([]byte(fdb_res), &vid_res)

    if err != nil {
        return nil, err
    }

    if len(vid_res) < end_pos {
        end_pos = len(vid_res)
    }

    if start_pos > end_pos {
        start_pos = end_pos
    }

    videos := &Videos{}
    for _, video_seo_title := range vid_res[start_pos:end_pos] {
        video := &Video{}

        video_raw, err := p.FDBVideos.Get(video_seo_title)

        if err != nil {
            continue
        }

        err = json.Unmarshal([]byte(video_raw), &video)

        if err != nil {
            continue
        }

        *videos = append(*videos, video)
    }

    return videos, nil
}

func (p Performer) DbGetCategories() (*Categories, error) {
    cat_seo_titles := p.FDBCategories.Keys

    categories := &Categories{}
    for _, cat_seo_title := range cat_seo_titles {
        category_raw, err := p.FDBCategories.Get(cat_seo_title)

        if err != nil {
            continue
        }

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
    fdb_res, err := p.FDBVideos.GetFilteredReverseRange(
        start_pos, end_pos, "(.*?)" + kword + "(.*?)")

    if err != nil {
        return nil, err
    }

    videos := &Videos{}
    for _, video_raw := range fdb_res {
        video := &Video{}

        err = json.Unmarshal([]byte(video_raw), &video)

        if err != nil {
            continue
        }

        *videos = append(*videos, video)
    }

    return videos.Range(start_pos, end_pos), nil
}

func (p Performer) DbAddVideo(video *Video) error {
    video_json, err := json.Marshal(video)

    err = p.FDBVideos.Set(video.SeoTitle, string(video_json))

    if err != nil {
        return errors.New("Video already exists")
    }

    for i:=0; i<video.Categories.Length(); i++ {
        err = p.DbAddCategory(video.Categories[i])

        if err != nil {
            return err
        }

        p.DbVideoToCategory(video.SeoTitle, video.Categories[i].SeoTitle)
    }

    return nil
}

func (p Performer) DbVideoToCategory(
  vid_seo_title string, cat_seo_title string) error {
    p.FDBCategoryVideos.Set(cat_seo_title, "[]")

    cat_data, err := p.FDBCategoryVideos.Get(cat_seo_title)

    if err != nil {
        return err
    }

    vid_res := []string{}
    err = json.Unmarshal([]byte(cat_data), &vid_res)

    if err != nil {
        return err
    }

    vid_res = append(vid_res, vid_seo_title)

    vid_list_json, err := json.Marshal(vid_res)

    if err != nil {
        return err
    }

    err = p.FDBCategoryVideos.Update(cat_seo_title, string(vid_list_json))

    if err != nil {
        return err
    }

    return nil
}

func (p Performer) DbAddCategory(category *Category) error {
    category_json, err := json.Marshal(category)

    if err != nil {
        return err
    }

    p.FDBCategories.Set(category.SeoTitle, string(category_json))

    return nil
}
