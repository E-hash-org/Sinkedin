package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sinkedin/models"
)

type CreateCommentInput struct {
	PostID    *uint    `json:"postId"`
	ParentID  *uint    `json:"parentId"`
	Content   string   `json:"content" binding:"required"`
	Type      string   `json:"type" binding:"required"`
	Tags      []string `json:"tags"`
	Hashtags  []string `json:"hashtags"`
}

func CreateComment(c *gin.Context) {
	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetUint("userId")
	comment := models.Comment{
		UserID:          userId,
		PostID:          input.PostID,
		ParentCommentID: input.ParentID,
		Type:            models.CommentType(input.Type),
		Content:         input.Content,
	}

	tx := models.DB.Begin()

	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Handle hashtags
	if len(input.Hashtags) > 0 {
		comment.ContainsHashtag = true
		for _, hashtagName := range input.Hashtags {
			var hashtag models.Hashtag
			if err := tx.Where("name = ?", hashtagName).FirstOrCreate(&hashtag, models.Hashtag{Name: hashtagName}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process hashtags"})
				return
			}
			if err := tx.Create(&models.CommentHashtag{CommentID: comment.ID, HashtagID: hashtag.ID}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate hashtags"})
				return
			}
		}
	}

	// Handle tags
	if len(input.Tags) > 0 {
		comment.ContainsTag = true
		for _, username := range input.Tags {
			var taggedUser models.User
			if err := tx.Where("username = ?", username).First(&taggedUser).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid tagged user"})
				return
			}
			if err := tx.Create(&models.CommentTag{CommentID: comment.ID, UserID: taggedUser.ID}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate tags"})
				return
			}
		}
	}

	// Update comment count on parent
	if comment.PostID != nil {
		tx.Model(&models.Post{}).Where("id = ?", *comment.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	}
	if comment.ParentCommentID != nil {
		tx.Model(&models.Comment{}).Where("id = ?", *comment.ParentCommentID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	}

	tx.Commit()

	// Load the complete comment with associations
	models.DB.Preload("User").Preload("Tags").Preload("Hashtags").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, comment)
}

func GetPostComments(c *gin.Context) {
	postId := c.Param("postId")

	var comments []models.Comment
	if err := models.DB.Preload("User").Preload("Tags").Preload("Hashtags").
		Where("post_id = ? AND parent_comment_id IS NULL", postId).
		Order("created_at desc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func GetComment(c *gin.Context) {
	id := c.Param("id")

	var comment models.Comment
	if err := models.DB.Preload("User").Preload("Tags").Preload("Hashtags").
		Preload("ParentComment").First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func UpdateComment(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetUint("userId")

	var comment models.Comment
	if err := models.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.Content = input.Content
	comment.Type = models.CommentType(input.Type)

	models.DB.Save(&comment)
	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetUint("userId")

	var comment models.Comment
	if err := models.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	tx := models.DB.Begin()

	// Update comment count on parent
	if comment.PostID != nil {
		tx.Model(&models.Post{}).Where("id = ?", *comment.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	}
	if comment.ParentCommentID != nil {
		tx.Model(&models.Comment{}).Where("id = ?", *comment.ParentCommentID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	}

	if err := tx.Delete(&comment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
