package handlers

import (	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sinkedin/models"
)

func ToggleLike(c *gin.Context) {
	likeType := models.LikeType(c.Param("type"))
	parentIdStr := c.Param("id")
	parentId64, err := strconv.ParseUint(parentIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	parentId := uint(parentId64)
	userId := c.GetUint("userId")

	if likeType != models.PostLike && likeType != models.CommentLike {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid like type"})
		return
	}

	var like models.Like
	result := models.DB.Where("user_id = ? AND parent_id = ? AND type = ?", userId, parentId, likeType).First(&like)

	tx := models.DB.Begin()

	if result.Error == gorm.ErrRecordNotFound {
		// Create new like
		like = models.Like{
			UserID:   userId,
			ParentID: uint(parentId),
			Type:     likeType,
		}

		if err := tx.Create(&like).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create like"})
			return
		}

		// Increment like count
		if likeType == models.PostLike {
			tx.Model(&models.Post{}).Where("id = ?", parentId).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
		} else {
			tx.Model(&models.Comment{}).Where("id = ?", parentId).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
		}

		tx.Commit()
		c.JSON(http.StatusCreated, gin.H{"message": "Liked successfully"})
	} else {
		// Remove existing like
		if err := tx.Delete(&like).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like"})
			return
		}

		// Decrement like count
		if likeType == models.PostLike {
			tx.Model(&models.Post{}).Where("id = ?", parentId).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
		} else {
			tx.Model(&models.Comment{}).Where("id = ?", parentId).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Unliked successfully"})
	}
}

func GetLikes(c *gin.Context) {
	likeType := models.LikeType(c.Param("type"))
	parentId := c.Param("id")

	var likes []models.Like
	if err := models.DB.Preload("User").
		Where("parent_id = ? AND type = ?", parentId, likeType).
		Order("created_at desc").Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch likes"})
		return
	}

	c.JSON(http.StatusOK, likes)
}
