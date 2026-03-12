package mysql

import (
	"MicroCraft/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostListFilter struct {
	Section string
	Status  string
	Sort    string
	WorkType string
	Offset  int
	Limit   int
}

type PostListItem struct {
	PostID    uint   `json:"post_id"`
	Section   string `json:"section"`
	Status    string `json:"status"`
	LikeCount int64  `json:"like_count"`
	ViewCount int64  `json:"view_count"`
	CreatedAt string `json:"created_at"`

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

func CountPosts(filter PostListFilter) (int64, error) {
	var total int64
	q := DB.Model(&model.Post{}).
		Where("section = ? AND status = ?", filter.Section, filter.Status)

	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func ListPosts(filter PostListFilter) ([]model.Post, error) {
	var list []model.Post
	q := DB.Model(&model.Post{}).
		Where("section = ? AND status = ?", filter.Section, filter.Status)
	switch filter.Sort {
	case "newest":
		q = q.Order("created_at DESC")
	case "view_count":
		q = q.Order("view_count DESC, like_count DESC, created_at DESC")
	default: // like_count
		q = q.Order("like_count DESC, view_count DESC, created_at DESC")
	}
	if err := q.Offset(filter.Offset).Limit(filter.Limit).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func GetPostByID(id uint) (*model.Post, error) {
	var p model.Post
	if err := DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func IncPostViewCount(id uint) error {
	return DB.Model(&model.Post{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func CountPostsWithWork(filter PostListFilter) (int64, error) {
	var total int64

	q := DB.Table("posts").
		Joins("JOIN works ON works.id = posts.work_id").
		Where("posts.section = ? AND posts.status = ?", filter.Section, filter.Status).
		Where("works.status = ?", "active")
	if filter.WorkType != "" {
		q = q.Where("works.type = ?", filter.WorkType)
	}
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func ListPostsWithWork(filter PostListFilter) ([]PostListItem, error) {
	type row struct {
		PostID    uint
		Section   string
		Status    string
		LikeCount int64
		ViewCount int64
		PostCreatedAt  string

		WorkID       uint
		UserID       uint
		Title        string
		Type         string
		Source       string
		FileURL      string
		ThumbnailURL string
		WorkStatus   string
		WorkCreatedAt string
	}

	var rows []row

	q := DB.Table("posts").
		Select(`
			posts.id as post_id,
			posts.section,
			posts.status,
			posts.like_count,
			posts.view_count,
			DATE_FORMAT(posts.created_at, '%Y-%m-%dT%H:%i:%s') as post_created_at,

			works.id as work_id,
			works.user_id,
			works.title,
			works.type,
			works.source,
			works.file_url,
			works.thumbnail_url,
			works.status as work_status,
			DATE_FORMAT(works.created_at, '%Y-%m-%dT%H:%i:%s') as work_created_at
		`).
		Joins("JOIN works ON works.id = posts.work_id").
		Where("posts.section = ? AND posts.status = ?", filter.Section, filter.Status).
		Where("works.status = ?", "active")

	if filter.WorkType != "" {
		q = q.Where("works.type = ?", filter.WorkType)
	}
	switch filter.Sort {
	case "newest":
		q = q.Order("posts.created_at DESC")
	case "view_count":
		q = q.Order("posts.view_count DESC, posts.like_count DESC, posts.created_at DESC")
	default: 
		q = q.Order("posts.like_count DESC, posts.view_count DESC, posts.created_at DESC")
	}

	if err := q.Offset(filter.Offset).Limit(filter.Limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]PostListItem, 0, len(rows))
	for _, r := range rows {
		var it PostListItem
		it.PostID = r.PostID
		it.Section = r.Section
		it.Status = r.Status
		it.LikeCount = r.LikeCount
		it.ViewCount = r.ViewCount
		it.CreatedAt = r.PostCreatedAt

		it.Work.WorkID = r.WorkID
		it.Work.UserID = r.UserID
		it.Work.Title = r.Title
		it.Work.Type = r.Type
		it.Work.Source = r.Source
		it.Work.FileURL = r.FileURL
		it.Work.ThumbnailURL = r.ThumbnailURL
		it.Work.Status = r.WorkStatus
		it.Work.CreatedAt = r.WorkCreatedAt

		list = append(list, it)
	}

	return list, nil
}

func UpdatePostStatus(postID uint, status string) error {
	return DB.Model(&model.Post{}).
		Where("id = ?", postID).
		Update("status", status).Error
}

func UnpublishPostsByWorkID(workID uint) error {
	return DB.Model(&model.Post{}).
		Where("work_id = ? AND status = ?", workID, "published").
		Update("status", "draft").Error
}

func GetPostByIDForUpdate(tx *gorm.DB, id uint) (*model.Post, error) {
	var p model.Post
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func IncPostLikeCountTx(tx *gorm.DB, postID uint, delta int64) error {
	return tx.Model(&model.Post{}).
		Where("id = ?", postID).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

type WorkWithAuthor struct {
	WorkID       uint   `json:"work_id"`
	UserID       uint   `json:"user_id"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	Source       string `json:"source"`
	FileURL      string `json:"file_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Status       string `json:"status"`
	AuthorNickname string `json:"author_nickname"`
}
func GetWorkWithAuthorByID(workID uint) (*WorkWithAuthor, error) {
	var out WorkWithAuthor
	err := DB.Table("works").
		Select(`
			works.id as work_id,
			works.user_id,
			works.title,
			works.type,
			works.source,
			works.file_url,
			works.thumbnail_url,
			works.status,
			users.nickname as author_nickname
		`).
		Joins("JOIN users ON users.id = works.user_id").
		Where("works.id = ? AND works.status = 'active'", workID).
		Scan(&out).Error
	if err != nil {
		return nil, err
	}
	return &out, nil
}
