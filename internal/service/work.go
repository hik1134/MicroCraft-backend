package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"MicroCraft/internal/dao/mysql"
	"MicroCraft/internal/model"
	perr "MicroCraft/pkg/errors"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

const (
	maxUploadSize = 8 << 20 
)

type UploadWorkResp struct {
	WorkID       uint   `json:"work_id"`
	FileURL      string `json:"file_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Source       string `json:"source"`
}

func UploadWorkLocal(ctx context.Context, userID uint, title, typ string, fileHeader *multipart.FileHeader) (*UploadWorkResp, error) {
	return uploadWork(ctx, userID, title, typ, fileHeader, "local_upload")
}

func UploadWorkPhoto(ctx context.Context, userID uint, title, typ string, fileHeader *multipart.FileHeader) (*UploadWorkResp, error) {
	return uploadWork(ctx, userID, title, typ, fileHeader, "photo_upload")
}

func uploadWork(ctx context.Context, userID uint, title, typ string, fileHeader *multipart.FileHeader, source string) (*UploadWorkResp, error) {
	if userID == 0 {
		return nil, perr.New(perr.AUTH_TOKEN_INVALID)
	}
	if fileHeader == nil {
		return nil, perr.New(perr.FILE_MISSING)
	}

	if fileHeader.Size <= 0 {
		return nil, perr.New(perr.FILE_MISSING)
	}
	if fileHeader.Size > maxUploadSize {
		return nil, perr.New(perr.FILE_TOO_LARGE)
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allow := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allow[ext] {
		return nil, perr.New(perr.FILE_TYPE_NOT_ALLOW)
	}

    section := ""
    switch source {
    case "local_upload":
        section = "exhibition"
    case "photo_upload":
        section = "life_texture"
    default:
        section = "exhibition" 
    }

    var w model.Work
    var fileURL string

    err := mysql.DB.Transaction(func(tx *gorm.DB) error {
        if typ == "" {
            typ = "order"
        }
        w = model.Work{
            UserID: userID,
            Title:  strings.TrimSpace(title),
            Type:   typ,
            Source: source,
            Status: "active",
            FileKey: "pending",
            FileURL: "pending",
        }
        if err := tx.Create(&w).Error; err != nil {
            return perr.Wrap(perr.WORK_CREATE_FAIL, err)
        }
        filename := uuid.New().String() + ext
        relDir := fmt.Sprintf("uploads/works/%d", w.ID)
        absDir := "./" + relDir
        if err := os.MkdirAll(absDir, 0755); err != nil {
            return perr.Wrap(perr.SAVE_FILE_FAIL, err)
        }

        relPath := fmt.Sprintf("%s/%s", relDir, filename)
        fileKey := relPath
        fileURL = "/" + relPath

        if err := tx.Model(&model.Work{}).Where("id = ?", w.ID).
            Updates(map[string]interface{}{
                "file_key":       fileKey,
                "file_url":       fileURL,
                "thumbnail_key":  fileKey,
                "thumbnail_url":  fileURL,
            }).Error; err != nil {
            return perr.Wrap(perr.DB_UPDATE_FAIL, err)
        }

        p := model.Post{
            WorkID:  w.ID,
            Section: section,        
            Status:  "published",   
            LikeCount: 0,
            ViewCount: 0,
        }
        if err := tx.Create(&p).Error; err != nil {
            return perr.Wrap(perr.DB_CREATE_FAIL, err)
        }

        return nil
    })
    if err != nil {
        return nil, err
    }

    return &UploadWorkResp{
        WorkID:       w.ID,
        FileURL:      fileURL,
        ThumbnailURL: fileURL,
        Source:       source,
    }, nil
}

type MyWorksQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Type     string `form:"type"`   
	Source   string `form:"source"` 
	Status   string `form:"status"` 
}

type MyWorksResp struct {
	List     []model.Work `json:"list"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
	Total    int64        `json:"total"`
}

func GetMyWorks(ctx context.Context, userID uint, q MyWorksQuery) (*MyWorksResp, error) {
	if userID == 0 {
		return nil, perr.New(perr.AUTH_TOKEN_INVALID)
	}

	page := q.Page
	if page <= 0 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize <= 0 {
		pageSize = 6
	}
	if pageSize > 50 {
		pageSize = 50
	}

	typ := strings.TrimSpace(q.Type)
	source := strings.TrimSpace(q.Source)
	status := strings.TrimSpace(q.Status)

	if status == "" {
		status = "active"
	}

	offset := (page - 1) * pageSize

	list, total, err := mysql.ListWorksByUser(userID, mysql.WorkListOpt{
		Type:   typ,
		Source: source,
		Status: status,
		Offset: offset,
		Limit:  pageSize,
	})
	if err != nil {
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	return &MyWorksResp{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

type WorkDetailResp struct {
	WorkID       uint   `json:"work_id"`
	UserID       uint   `json:"user_id"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	Source       string `json:"source"`
	FileURL      string `json:"file_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"` // 也可以 time.Time，看你前端偏好
}

func GetWorkDetail(userID uint, workID uint) (*model.Work, error) {
	if userID == 0 || workID == 0 {
		return nil, perr.New(perr.INVALID_PARAM)
	}

	w, err := mysql.GetWorkByID(workID)
	if err != nil {
		if mysql.IsNotFound(err) {
			return nil, perr.New(perr.WORK_NOT_FOUND)
		}
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	if w.UserID != userID {
		return nil, perr.New(perr.WORK_FORBIDDEN)
	}

	if strings.ToLower(w.Status) != "active" {
		return nil, perr.New(perr.WORK_NOT_FOUND)
	}

	return w, nil
}

func DeleteWork(userID uint, workID uint) error {
	if userID == 0 || workID == 0 {
		return perr.New(perr.INVALID_PARAM)
	}

	w, err := mysql.GetWorkByID(workID)
	if err != nil {
		if mysql.IsNotFound(err) {
			return perr.New(perr.WORK_NOT_FOUND)
		}
		return perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	if w.UserID != userID {
		return perr.New(perr.WORK_FORBIDDEN)
	}

	if strings.ToLower(w.Status) != "active" {
		return perr.New(perr.WORK_NOT_FOUND)
	}

	if err := mysql.SoftDeleteWork(workID); err != nil {
		return perr.Wrap(perr.DB_UPDATE_FAIL, err)
	}

	if err := mysql.UnpublishPostsByWorkID(workID); err != nil {
		return perr.Wrap(perr.DB_UPDATE_FAIL, err)
	}

	return nil
}