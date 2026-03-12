package service

import (
	"context"
	"math"
	"strings"

	"MicroCraft/internal/dao/mysql"
	perr "MicroCraft/pkg/errors"
)

type ExhibitionListReq struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Sort     string `form:"sort"` 
}

type ExhibitionListResp struct {
	List     []mysql.PostListItem `json:"list"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
	Total    int64                `json:"total"`
}

func GetExhibitionPosts(ctx context.Context, req ExhibitionListReq) (*ExhibitionListResp, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 6
	}
	pageSize = int(math.Min(float64(pageSize), 50))

	sort := strings.TrimSpace(req.Sort)
	if sort == "" {
		sort = "like_count"
	}

	filter := mysql.PostListFilter{
		Section: "exhibition",
		Status:  "published",
		Sort:    sort,
		Offset:  (page - 1) * pageSize,
		Limit:   pageSize,
	}

	total, err := mysql.CountPostsWithWork(filter)
	if err != nil {
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	list, err := mysql.ListPostsWithWork(filter)
	if err != nil {
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	return &ExhibitionListResp{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

type LifeTextureListReq struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Type     string `form:"type"` 
	Sort     string `form:"sort"` 
}

type LifeTextureListResp struct {
	List     []mysql.PostListItem `json:"list"` 
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
	Total    int64                `json:"total"`
}

func GetLifeTexturePosts(ctx context.Context, req LifeTextureListReq) (*LifeTextureListResp, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 6
	}
	pageSize = int(math.Min(float64(pageSize), 50))

	mode := strings.TrimSpace(req.Type)
	if mode == "" {
		mode = "new"
	}
	if mode != "new" && mode != "hot" {
		return nil, perr.New(perr.INVALID_PARAM)
	}

	workType := strings.TrimSpace(req.Sort)
	if workType != "" {
		allow := map[string]bool{"order": true, "logic": true, "time": true, "flow": true}
		if !allow[workType] {
			return nil, perr.New(perr.INVALID_PARAM)
		}
	}

	order := "newest"
	if mode == "hot" {
		order = "like_count"
	}

	filter := mysql.PostListFilter{
		Section:  "life_texture",
		Status:   "published",
		Sort:     order,
		WorkType: workType, 
		Offset:   (page - 1) * pageSize,
		Limit:    pageSize,
	}

	total, err := mysql.CountPostsWithWork(filter)
	if err != nil {
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	list, err := mysql.ListPostsWithWork(filter)
	if err != nil {
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	return &LifeTextureListResp{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

type AuthorInfo struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
}