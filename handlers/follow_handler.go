package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sinkedin/models"
)

func ToggleFollow(c *gin.Context) {
	username := c.Param("username")
	followerId := c.GetUint("userId")

	var targetUser models.User
	if err := models.DB.Where("username = ?", username).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if targetUser.ID == followerId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	var follow models.Follow
	result := models.DB.Where("follower_id = ? AND following_id = ?", followerId, targetUser.ID).First(&follow)

	tx := models.DB.Begin()

	if result.Error == gorm.ErrRecordNotFound {
		// Create new follow
		follow = models.Follow{
			FollowerID:  followerId,
			FollowingID: targetUser.ID,
		}

		if err := tx.Create(&follow).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
			return
		}

		// Update follower/following counts
		tx.Model(&models.User{}).Where("id = ?", followerId).UpdateColumn("following_count", gorm.Expr("following_count + ?", 1))
		tx.Model(&models.User{}).Where("id = ?", targetUser.ID).UpdateColumn("followers_count", gorm.Expr("followers_count + ?", 1))

		tx.Commit()
		c.JSON(http.StatusCreated, gin.H{"message": "Following successfully"})
	} else {
		// Remove existing follow
		if err := tx.Delete(&follow).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
			return
		}

		// Update follower/following counts
		tx.Model(&models.User{}).Where("id = ?", followerId).UpdateColumn("following_count", gorm.Expr("following_count - ?", 1))
		tx.Model(&models.User{}).Where("id = ?", targetUser.ID).UpdateColumn("followers_count", gorm.Expr("followers_count - ?", 1))

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
	}
}

func GetFollowers(c *gin.Context) {
	username := c.Param("username")

	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var followers []models.User
	if err := models.DB.Table("users").
		Select("users.*").
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", user.ID).
		Order("follows.created_at desc").
		Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	c.JSON(http.StatusOK, followers)
}

func GetFollowing(c *gin.Context) {
	username := c.Param("username")

	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var following []models.User
	if err := models.DB.Table("users").
		Select("users.*").
		Joins("JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ?", user.ID).
		Order("follows.created_at desc").
		Find(&following).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch following"})
		return
	}

	c.JSON(http.StatusOK, following)
}
