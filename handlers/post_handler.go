package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sinkedin/models"
)

type CreatePostInput struct {
	Content    string   `json:"content" binding:"required"`
	ImageURL   string   `json:"imageURL"`
	Tags       []string `json:"tags"`
	Hashtags   []string `json:"hashtags"`
	IsQuote    bool     `json:"isQuote"`
	QuoteLines string   `json:"quoteLines"`
}

func CreatePost(c *gin.Context) {
	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetUint("userId")
	post := models.Post{
		UserID:     userId,
		Content:    input.Content,
		ImageURL:   input.ImageURL,
		HasImage:   input.ImageURL != "",
		IsQuote:    input.IsQuote,
		QuoteLines: input.QuoteLines,
	}

	// Start a transaction
	tx := models.DB.Begin()

	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Handle hashtags
	if len(input.Hashtags) > 0 {
		post.HasHashtag = true
		for _, hashtagName := range input.Hashtags {
			var hashtag models.Hashtag
			// Find or create hashtag
			if err := tx.Where("name = ?", hashtagName).FirstOrCreate(&hashtag, models.Hashtag{Name: hashtagName}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process hashtags"})
				return
			}
			// Increment counter
			tx.Model(&hashtag).UpdateColumn("counter", hashtag.Counter+1)
			// Create association
			if err := tx.Create(&models.PostHashtag{PostID: post.ID, HashtagID: hashtag.ID}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate hashtags"})
				return
			}
		}
	}

	// Handle tags
	if len(input.Tags) > 0 {
		post.HasTag = true
		for _, username := range input.Tags {
			var taggedUser models.User
			if err := tx.Where("username = ?", username).First(&taggedUser).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid tagged user"})
				return
			}
			if err := tx.Create(&models.PostTag{PostID: post.ID, UserID: taggedUser.ID}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate tags"})
				return
			}
		}
	}

	tx.Commit()

	// Load the complete post with associations
	models.DB.Preload("User").Preload("Tags").Preload("Hashtags").First(&post, post.ID)

	c.JSON(http.StatusCreated, post)
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := models.DB.Preload("User").Preload("Tags").Preload("Hashtags").Order("created_at desc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := models.DB.Preload("User").Preload("Tags").Preload("Hashtags").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetUint("userId")

	var post models.Post
	if err := models.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Content = input.Content
	post.ImageURL = input.ImageURL
	post.HasImage = input.ImageURL != ""
	post.IsQuote = input.IsQuote
	post.QuoteLines = input.QuoteLines

	models.DB.Save(&post)
	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetUint("userId")

	var post models.Post
	if err := models.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	models.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
