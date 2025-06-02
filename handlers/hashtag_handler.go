package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sinkedin/models"
)

func GetTrendingHashtags(c *gin.Context) {
	var hashtags []models.Hashtag
	if err := models.DB.Order("counter desc").Limit(10).Find(&hashtags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trending hashtags"})
		return
	}

	c.JSON(http.StatusOK, hashtags)
}

func GetHashtagPosts(c *gin.Context) {
	hashtagName := c.Param("name")

	var hashtag models.Hashtag
	if err := models.DB.Where("name = ?", hashtagName).First(&hashtag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hashtag not found"})
		return
	}

	var posts []models.Post
	if err := models.DB.Preload("User").Preload("Tags").Preload("Hashtags").
		Joins("JOIN post_hashtags ON posts.id = post_hashtags.post_id").
		Where("post_hashtags.hashtag_id = ?", hashtag.ID).
		Order("posts.created_at desc").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
