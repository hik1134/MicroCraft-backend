package service

import (
	"context"
	"errors"

	"MicroCraft/internal/dao/mysql"
	perr "MicroCraft/pkg/errors"
	"gorm.io/gorm"
)

type PostDetailResp struct {
	PostID    uint   `json:"post_id"`
	Section   string `json:"section"`
	Status    string `json:"status"`
	LikeCount int64  `json:"like_count"`
	ViewCount int64  `json:"view_count"`
	CreatedAt string `json:"created_at"`

	Liked bool `json:"liked"` 
	Author AuthorInfo `json:"author"`

	Work struct {
		WorkID       uint   `json:"work_id"`
		UserID       uint   `json:"user_id"`
		Title        string `json:"title"`
		Type         string `json:"type"`
		Source       string `json:"source"`
		FileURL      string `json:"file_url"`
		ThumbnailURL string `json:"thumbnail_url"`
		Status       string `json:"status"`
		CreatedAt    string `json:"created_at"`
	} `json:"work"`
}

func GetPostDetail(ctx context.Context, postID uint, userID uint) (*PostDetailResp, error) {
	p, err := mysql.GetPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, perr.New(perr.POST_NOT_FOUND)
		}
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	w, err := mysql.GetWorkByID(p.WorkID)
	if err != nil {
		if mysql.IsNotFound(err) {
			return nil, perr.New(perr.POST_NOT_FOUND) // work 没了，这帖子也当不存在
		}
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	isOwner := (userID != 0 && w.UserID == userID)

	if p.Status == "draft" && !isOwner {
		return nil, perr.New(perr.POST_FORBIDDEN)
	}

	if p.Status == "published" {
		if err := mysql.IncPostViewCount(p.ID); err != nil {
			return nil, perr.Wrap(perr.DB_UPDATE_FAIL, err)
		}
		p.ViewCount++
	}

	liked := false
	if userID != 0 {
		liked, err = mysql.IsPostLiked(p.ID, userID)
		if err != nil {
			return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
		}
	}

	resp := &PostDetailResp{
		PostID:    p.ID,
		Section:   p.Section,
		Status:    p.Status,
		LikeCount: p.LikeCount,
		ViewCount: p.ViewCount,
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
		Liked:     liked,
	}

	resp.Work.WorkID = w.ID
	resp.Work.UserID = w.UserID
	resp.Work.Title = w.Title
	resp.Work.Type = w.Type
	resp.Work.Source = w.Source
	resp.Work.FileURL = w.FileURL
	resp.Work.ThumbnailURL = w.ThumbnailURL
	resp.Work.Status = w.Status
	resp.Work.CreatedAt = w.CreatedAt.Format("2006-01-02T15:04:05-07:00")

	return resp, nil
}

func GetLifeTexturePostDetail(ctx context.Context, postID uint, userID uint) (*PostDetailResp, error) {
    resp, err := GetPostDetail(ctx, postID, userID)
    if err != nil {
        return nil, err
    }

    if resp.Section != "life_texture" {
        return nil, perr.New(perr.POST_NOT_FOUND)
    }

    return resp, nil
}