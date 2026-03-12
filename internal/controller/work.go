package controller

import (
	"fmt"
	"path/filepath"
	"strings"
	"strconv"
	"MicroCraft/internal/middleware"
	"MicroCraft/internal/service"
	perr "MicroCraft/pkg/errors"
	"MicroCraft/pkg/response"
	"github.com/gin-gonic/gin"
)

func UploadWorkLocalHandler(c *gin.Context) {
	v, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_MISSING))
		return
	}
	userID, ok := v.(uint)
	if !ok || userID == 0 {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_INVALID))
		return
	}
	title := strings.TrimSpace(c.PostForm("title"))
	typ := strings.TrimSpace(c.PostForm("type"))

	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		response.FailByErr(c, perr.New(perr.FILE_MISSING))
		return
	}
	resp, err := service.UploadWorkLocal(c.Request.Context(), userID, title, typ, fileHeader)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	relPath := strings.TrimPrefix(resp.FileURL, "/")
	savePath := filepath.FromSlash("./" + relPath)
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		response.FailByErr(c, perr.Wrap(perr.SAVE_FILE_FAIL, fmt.Errorf("save file: %w", err)))
		return
	}
	response.OK(c, resp)
}

func UploadWorkPhotoHandler(c *gin.Context) {
	v, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_MISSING))
		return
	}
	userID, ok := v.(uint)
	if !ok || userID == 0 {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_INVALID))
		return
	}
	title := strings.TrimSpace(c.PostForm("title"))
	typ := strings.TrimSpace(c.PostForm("type"))
	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		response.FailByErr(c, perr.New(perr.FILE_MISSING))
		return
	}
	resp, err := service.UploadWorkPhoto(c.Request.Context(), userID, title, typ, fileHeader)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	relPath := strings.TrimPrefix(resp.FileURL, "/")
	savePath := filepath.FromSlash("./" + relPath)
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		response.FailByErr(c, perr.Wrap(perr.SAVE_FILE_FAIL, fmt.Errorf("save file: %w", err)))
		return
	}
	response.OK(c, resp)
}

func GetMyWorksHandler(c *gin.Context) {
	v, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_MISSING))
		return
	}
	userID, ok := v.(uint)
	if !ok || userID == 0 {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_INVALID))
		return
	}
	var q service.MyWorksQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	resp, err := service.GetMyWorks(c.Request.Context(), userID, q)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, resp)
}

func GetWorkDetailHandler(c *gin.Context) {
	v, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_MISSING))
		return
	}
	userID, ok := v.(uint)
	if !ok || userID == 0 {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_INVALID))
		return
	}
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	workID := uint(id64)
	w, err := service.GetWorkDetail(userID, workID)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, w)
}

func DeleteWorkHandler(c *gin.Context) {
	v, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_MISSING))
		return
	}
	userID, ok := v.(uint)
	if !ok || userID == 0 {
		response.FailByErr(c, perr.New(perr.AUTH_TOKEN_INVALID))
		return
	}
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	workID := uint(id64)
	if err := service.DeleteWork(userID, workID); err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, gin.H{
		"work_id": workID,
		"status":  "deleted",
	})
}