package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shrey209/Like-Service/model"
	"github.com/shrey209/Like-Service/service"
)

type LikeController struct {
	Service *service.LikeService
}

func NewLikeController(service *service.LikeService) *LikeController {
	return &LikeController{Service: service}
}

func (lc *LikeController) RegisterRoutes(router *gin.Engine) {
	router.POST("/like/:postID", lc.LikePost)
	router.GET("/like/:postID", lc.GetLikeCount)
	router.POST("/likes/batch", lc.BatchLikePosts)
	router.POST("/like/init/:postID", lc.InitPostLike)
}

func (lc *LikeController) LikePost(c *gin.Context) {
	postID := c.Param("postID")
	err := lc.Service.LikePost(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post liked"})
}

func (lc *LikeController) GetLikeCount(c *gin.Context) {
	postID := c.Param("postID")
	count, err := lc.Service.GetPostLikes(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"postID": postID, "likes": count})
}

func (lc *LikeController) BatchLikePosts(c *gin.Context) {
	var likes []model.PostLike
	if err := c.ShouldBindJSON(&likes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := lc.Service.BatchLikePosts(likes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Batch update successful"})
}

func (lc *LikeController) InitPostLike(c *gin.Context) {
	postID := c.Param("postID")
	err := lc.Service.InitPostLike(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post initialized"})
}
