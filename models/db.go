package models

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

type User struct {
    ID            uint           `gorm:"primaryKey;type:serial" json:"id"`
    Name          string         `gorm:"type:varchar(100);not null" json:"name"`
    Username      string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
    Email         string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"email"`
    PhoneNumber   string         `gorm:"type:varchar(20)" json:"phoneNumber"`
    Password      string         `gorm:"type:varchar(255);not null" json:"-"` 
    Bio           string         `gorm:"type:varchar(500)" json:"bio"`
    FollowersCount int           `gorm:"default:0" json:"followersCount"`
    FollowingCount int           `gorm:"default:0" json:"followingCount"`
    DOB           *time.Time     `json:"dob"`
    PhotoURL      string         `gorm:"type:varchar(255)" json:"photoURL"`
    BannerURL     string         `gorm:"type:varchar(255)" json:"bannerURL"`
    CreatedAt     time.Time      `json:"createdAt"`
    UpdatedAt     time.Time      `json:"updatedAt"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Post struct {
    ID           uint           `gorm:"primaryKey;type:serial" json:"id"`
    UserID       uint           `gorm:"not null;index" json:"userId"`
    User         User           `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
    Content      string         `gorm:"type:text;not null" json:"content"` 
    HasImage     bool           `gorm:"default:false" json:"hasImage"`
    HasTag       bool           `gorm:"default:false" json:"hasTag"`
    HasHashtag   bool           `gorm:"default:false" json:"hasHashtag"`
    ImageURL     string         `gorm:"type:varchar(255)" json:"imageURL"`
    LikeCount    int            `gorm:"default:0" json:"likeCount"`
    CommentCount int            `gorm:"default:0" json:"commentCount"`
    IsQuote      bool           `gorm:"default:false" json:"isQuote"`
    QuoteLines   string         `gorm:"type:varchar(500)" json:"quoteLines"`
    CreatedAt    time.Time      `gorm:"index" json:"createdAt"` 
    UpdatedAt    time.Time      `json:"updatedAt"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Many-to-many relationships
    Hashtags     []Hashtag      `gorm:"many2many:post_hashtags;" json:"hashtags"`
    Tags         []User         `gorm:"many2many:post_tags;" json:"tags"`
}

type Hashtag struct {
    ID        uint           `gorm:"primaryKey;type:serial" json:"id"`
    Name      string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
    Counter   int            `gorm:"default:0" json:"counter"` 
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CommentType string

const (
    NormalComment CommentType = "normal"
    GifComment    CommentType = "gif"
)

// type Comment struct {
//     ID             uint           `gorm:"primaryKey;type:serial" json:"id"`
//     UserID         uint           `gorm:"not null;index" json:"userId"`
//     User           User           `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
//     Type           CommentType    `gorm:"type:varchar(10);not null;default:'normal'" json:"type"`
//     Content        string         `gorm:"type:varchar(500)" json:"content"`
//     ContainsTag    bool           `gorm:"default:false" json:"containsTag"`
//     ContainsHashtag bool          `gorm:"default:false" json:"containsHashtag"`
//     IsPostComment  bool           `gorm:"default:true;index" json:"isPostComment"` 
//     ParentID       uint           `gorm:"not null;index" json:"parentId"` 
//     LikeCount      int            `gorm:"default:0" json:"likeCount"`
//     CommentCount   int            `gorm:"default:0" json:"commentCount"`
//     CreatedAt      time.Time      `gorm:"index" json:"createdAt"` 
//     UpdatedAt      time.Time      `json:"updatedAt"`
//     DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
    
//     Tags           []User         `gorm:"many2many:comment_tags;" json:"tags"`
//     Hashtags       []Hashtag      `gorm:"many2many:comment_hashtags;" json:"hashtags"`
// }
// Structure and gorm support issue in the above struct <--- but plz dont delete anyone cuz it is useful for reference 

type Comment struct {
    ID              uint           `gorm:"primaryKey;type:serial" json:"id"`
    UserID          uint           `gorm:"not null;index" json:"userId"`
    User            User           `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
    PostID          *uint          `gorm:"index" json:"postId"` 
    Post            *Post          `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE" json:"post"`
    ParentCommentID *uint          `gorm:"index" json:"parentCommentId"`
    ParentComment   *Comment       `gorm:"foreignKey:ParentCommentID;references:ID;constraint:OnDelete:CASCADE" json:"parentComment"`
    Type            CommentType    `gorm:"type:varchar(10);not null;default:'normal'" json:"type"`
    Content         string         `gorm:"type:varchar(500)" json:"content"`
    ContainsTag     bool           `gorm:"default:false" json:"containsTag"`
    ContainsHashtag bool           `gorm:"default:false" json:"containsHashtag"`
    LikeCount       int            `gorm:"default:0" json:"likeCount"`
    CommentCount    int            `gorm:"default:0" json:"commentCount"`
    CreatedAt       time.Time      `gorm:"index" json:"createdAt"`
    UpdatedAt       time.Time      `json:"updatedAt"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
    Tags            []User         `gorm:"many2many:comment_tags;" json:"tags"`
    Hashtags        []Hashtag      `gorm:"many2many:comment_hashtags;" json:"hashtags"`
}

type LikeType string

const (
    PostLike    LikeType = "post"
    CommentLike LikeType = "comment"
)

type Like struct {
    ID        uint           `gorm:"primaryKey;type:serial" json:"id"`
    UserID    uint           `gorm:"not null;uniqueIndex:idx_user_parent_type" json:"userId"`
    User      User           `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
    ParentID  uint           `gorm:"not null;uniqueIndex:idx_user_parent_type" json:"parentId"`
    Type      LikeType       `gorm:"type:varchar(10);not null;uniqueIndex:idx_user_parent_type" json:"type"`
    CreatedAt time.Time      `json:"createdAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Follow struct {
    FollowerID  uint           `gorm:"primaryKey" json:"followerId"`
    FollowingID uint           `gorm:"primaryKey" json:"followingId"`
    Follower    User           `gorm:"foreignKey:FollowerID;references:ID;constraint:OnDelete:CASCADE" json:"follower"`
    Following   User           `gorm:"foreignKey:FollowingID;references:ID;constraint:OnDelete:CASCADE" json:"following"`
    CreatedAt   time.Time      `json:"createdAt"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type PostHashtag struct {
    PostID    uint      `gorm:"primaryKey;index:idx_post_hashtag"`
    HashtagID uint      `gorm:"primaryKey;index:idx_post_hashtag"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type PostTag struct {
    PostID    uint      `gorm:"primaryKey;index:idx_post_tag"`
    UserID    uint      `gorm:"primaryKey;index:idx_post_tag"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type CommentTag struct {
    CommentID uint      `gorm:"primaryKey;index:idx_comment_tag"`
    UserID    uint      `gorm:"primaryKey;index:idx_comment_tag"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type CommentHashtag struct {
    CommentID uint      `gorm:"primaryKey;index:idx_comment_hashtag"`
    HashtagID uint      `gorm:"primaryKey;index:idx_comment_hashtag"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func SetupDB() {
    host := os.Getenv("DB_HOST")
    if host == "" {
        host = "localhost"
    }
    
    user := os.Getenv("DB_USER")
    if user == "" {
        user = "postgres"
    }
    
    password := os.Getenv("DB_PASSWORD")
    if password == "" {
        password = "postgres"
    }
    
    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        dbName = "sinkedin"
    }
    
    port := os.Getenv("DB_PORT")
    if port == "" {
        port = "5432"
    }

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", 
        host, user, password, dbName, port)
    
    dbLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), 
        logger.Config{
            SlowThreshold:             time.Second, 
            LogLevel:                  logger.Info, 
            IgnoreRecordNotFoundError: true,        
            Colorful:                  true,        
        },
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: dbLogger,
    })

    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        log.Fatalf("Failed to get DB instance: %v", err)
    }
    
    sqlDB.SetMaxIdleConns(10)
    
    sqlDB.SetMaxOpenConns(100)
      sqlDB.SetConnMaxLifetime(time.Hour)

    log.Println("Setting up database connection...")
    DB = db
    log.Println("Database connection established successfully")
}