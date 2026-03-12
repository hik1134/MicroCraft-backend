package controller

import (
	"strings"
	"strconv"
	"MicroCraft/internal/middleware"
	"MicroCraft/internal/config"
	"MicroCraft/internal/service"
	perr "MicroCraft/pkg/errors"
	"MicroCraft/pkg/response"
	"MicroCraft/pkg/utils"
	"github.com/gin-gonic/gin"
)

func parseUintParam(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func GetExhibitionPostsHandler(c *gin.Context) {
	var req service.ExhibitionListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	resp, err := service.GetExhibitionPosts(c.Request.Context(), req)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, resp)
}

func GetPostDetailHandler(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	postID, err := parseUintParam(idStr) // 你可以自己写个小工具；下面我给一个
	if err != nil || postID == 0 {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	var userID uint = 0
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	if auth != "" && config.Conf != nil && config.Conf.Jwt.Secret != "" {
		tokenStr := auth
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			tokenStr = strings.TrimSpace(auth[7:])
		}
		if tokenStr != "" {
			uid, e := utils.ParseToken(config.Conf.Jwt.Secret, tokenStr)
			if e == nil {
				userID = uid 
			}
		}
	}
	resp, err := service.GetPostDetail(c.Request.Context(), uint(postID), userID)
	if err != nil {
		response.FailByErr(c, err)
		return
	}

	response.OK(c, resp)
}

func GetLifeTexturePostsHandler(c *gin.Context) {
    var req service.LifeTextureListReq
    if err := c.ShouldBindQuery(&req); err != nil {
        response.FailByErr(c, perr.New(perr.INVALID_PARAM))
        return
    }
    req.Type = strings.TrimSpace(req.Type)
    req.Sort = strings.TrimSpace(req.Sort)
    resp, err := service.GetLifeTexturePosts(c.Request.Context(), req)
    if err != nil {
        response.FailByErr(c, err)
        return
    }
    response.OK(c, resp)
}

func GetLifeTexturePostDetailHandler(c *gin.Context) {
    idStr := strings.TrimSpace(c.Param("id"))
    id64, err := parseUintParam(idStr)
    if err != nil || id64 == 0 {
        response.FailByErr(c, perr.New(perr.INVALID_PARAM))
        return
    }
    postID := uint(id64)
    var userID uint = 0
    auth := strings.TrimSpace(c.GetHeader("Authorization"))
    if auth != "" && config.Conf != nil && config.Conf.Jwt.Secret != "" {
        tokenStr := auth
        if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
            tokenStr = strings.TrimSpace(auth[7:])
        }
        if tokenStr != "" {
            uid, e := utils.ParseToken(config.Conf.Jwt.Secret, tokenStr)
            if e == nil {
                userID = uid
            }
        }
    }
    resp, err := service.GetLifeTexturePostDetail(c.Request.Context(), postID, userID)
    if err != nil {
        response.FailByErr(c, err)
        return
    }

    response.OK(c, resp)
}

func UnpublishPostHandler(c *gin.Context) {
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
	idStr := strings.TrimSpace(c.Param("post_id"))
	id64, err := parseUintParam(idStr)
	if err != nil || id64 == 0 {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	postID := uint(id64)
	if err := service.UnpublishPost(c.Request.Context(), userID, postID); err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, gin.H{
		"post_id": postID,
		"status":  "draft",
	})
}

func ToggleLikePostHandler(c *gin.Context) {
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
	idStr := strings.TrimSpace(c.Param("post_id"))
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		response.FailByErr(c, perr.New(perr.INVALID_PARAM))
		return
	}
	postID := uint(id64)
	resp, err := service.ToggleLikePost(c.Request.Context(), userID, postID)
	if err != nil {
		response.FailByErr(c, err)
		return
	}
	response.OK(c, resp)
}