package router

import (
	"MicroCraft/internal/controller"
	"MicroCraft/internal/middleware"
	"path/filepath"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
    	AllowAllOrigins: true,
    	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    	AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
    	ExposeHeaders: []string{"Content-Length"},
    	AllowCredentials: false,
    	MaxAge: 12 * time.Hour,
	}))
	r.MaxMultipartMemory = 8 << 20
	uploadDir, _ := filepath.Abs("./uploads")
	r.Static("/uploads", uploadDir)

	auth := r.Group("/auth")
	{
		auth.POST("/email/code", controller.SendEmailCodeHandler)
		auth.POST("/register", controller.RegisterHandler)
		auth.POST("/login", controller.LoginHandler)
		auth.GET("/me", middleware.AuthJWT(), controller.MeHandler)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthJWT())
	{
		works := api.Group("/works")
		{
			works.POST("/upload/local", controller.UploadWorkLocalHandler)
			works.POST("/upload/photo", controller.UploadWorkPhotoHandler)
			works.GET("/mine", controller.GetMyWorksHandler)
			works.GET("/:id", controller.GetWorkDetailHandler)
			works.DELETE("/:id", controller.DeleteWorkHandler)
		}
	}

	posts := r.Group("/posts")
	{
		posts.GET("/exhibition", controller.GetExhibitionPostsHandler)
		posts.GET("/:id", controller.GetPostDetailHandler)
		posts.GET("/life-texture", controller.GetLifeTexturePostsHandler)
		posts.GET("/life-texture/:id", controller.GetPostDetailHandler)
	}

	postsOwner := r.Group("/posts")
	postsOwner.Use(middleware.AuthJWT())
	{
		postsOwner.POST("/:post_id/unpublish", controller.UnpublishPostHandler)
		postsOwner.POST("/:post_id/like", controller.ToggleLikePostHandler)
	}

	r.GET("/carriers", controller.GetCarriersHandler)

	return r
}