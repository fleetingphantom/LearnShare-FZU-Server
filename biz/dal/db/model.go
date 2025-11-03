package db

import (
	"LearnShare/biz/model/module"
	"time"
)

type User struct {
	UserID          int64     `json:"user_id" db:"user_id"`
	Username        string    `json:"username" db:"username"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	Email           string    `json:"email" db:"email"`
	CollegeID       *int64    `json:"college_id,omitempty" db:"college_id"`
	MajorID         *int64    `json:"major_id,omitempty" db:"major_id"`
	AvatarURL       *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	ReputationScore int64     `json:"reputation_score" db:"reputation_score"`
	RoleID          int64     `json:"role_id" db:"role_id"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func (u User) ToUserModule() *module.User {
	user := &module.User{
		UserId:          u.UserID,
		Username:        u.Username,
		Email:           u.Email,
		ReputationScore: u.ReputationScore,
		RoleId:          u.RoleID,
		Status:          u.Status,
	}

	if u.AvatarURL != nil {
		user.AvatarUrl = *u.AvatarURL
	}
	if u.CollegeID != nil {
		user.CollegeId = *u.CollegeID
	}
	if u.MajorID != nil {
		user.MajorId = *u.MajorID
	}
	return user
}

// ToResourceModule 将db.Resource转换为model.Resource
func (r Resource) ToResourceModule() *module.Resource {
	var tags []*module.ResourceTag
	for _, t := range r.Tags {
		tags = append(tags, t.ToResourceTagModule())
	}

	return &module.Resource{
		ResourceId:    r.ResourceID,
		Title:         r.ResourceName,
		Description:   &r.Description,
		FilePath:      r.FilePath,
		FileType:      r.FileType,
		FileSize:      r.FileSize,
		UploaderId:    r.UploaderID,
		CourseId:      r.CourseID,
		DownloadCount: r.DownloadCount,
		AverageRating: r.AverageRating,
		RatingCount:   r.RatingCount,
		Status:        convertStatus(r.Status),
		CreatedAt:     r.CreatedAt.Unix(),
		Tags:          tags,
	}
}

func convertStatus(status string) int32 {
	switch status {
	case "pending_review":
		return 1
	case "normal":
		return 0
	case "low_quality":
		return 2
	default:
		return 1 // 默认为待审核
	}
}

// ToResourceTagModule 将db.ResourceTag转换为model.ResourceTag
func (t ResourceTag) ToResourceTagModule() *module.ResourceTag {
	return &module.ResourceTag{
		TagId:   t.TagID,
		TagName: t.TagName,
	}
}

// ToResourceCommentModule 将db.ResourceComment转换为model.ResourceComment
func (c ResourceComment) ToResourceCommentModule() *module.ResourceComment {
	var parentId int64
	if c.ParentID != nil {
		parentId = *c.ParentID
	}

	var status module.ResourceCommentStatus
	status, _ = module.ResourceCommentStatusFromString(c.Status)

	return &module.ResourceComment{
		CommentId:  c.CommentID,
		UserId:     c.UserID,
		ResourceId: c.ResourceID,
		Content:    c.Content,
		ParentId:   parentId,
		Likes:      c.Likes,
		IsVisible:  c.IsVisible,
		Status:     status,
		CreatedAt:  c.CreatedAt.Unix(),
	}
}

// ToResourceRatingModule 将db.ResourceRating转换为model.ResourceRating
func (r ResourceRating) ToResourceRatingModule() *module.ResourceRating {
	return &module.ResourceRating{
		RatingId:       r.RatingID,
		UserId:         r.UserID,
		ResourceId:     r.ResourceID,
		Recommendation: r.Recommendation * 10, // 转换为0-50的浮点数
		IsVisible:      r.IsVisible,
		CreatedAt:      r.CreatedAt.Unix(),
	}
}

type Resource struct {
	ResourceID    int64         `gorm:"primaryKey;autoIncrement"`
	ResourceName  string        `gorm:"column:resource_name;size:255;not null"`
	Description   string        `gorm:"type:text"`
	FilePath      string        `gorm:"column:resource_url;size:255;not null"`
	FileType      string        `gorm:"column:type;size:50;not null"`
	FileSize      int64         `gorm:"column:size;not null"`
	UploaderID    int64         `gorm:"not null"`
	CourseID      int64         `gorm:"not null"`
	DownloadCount int64         `gorm:"default:0"`
	AverageRating float64       `gorm:"default:0.0"`
	RatingCount   int64         `gorm:"default:0"`
	Status        string        `gorm:"type:enum('normal','low_quality','pending_review');default:'pending_review'"`
	CreatedAt     time.Time     `gorm:"autoCreateTime"`
	Tags          []ResourceTag `gorm:"many2many:resource_tags;joinForeignKey:resource_id;joinReferences:tag_id"`
}

type ResourceTag struct {
	TagID   int64  `gorm:"primaryKey;autoIncrement;table:tags"`
	TagName string `gorm:"size:50;unique;not null"`
}

type ResourceTagMapping struct {
	ResourceID int64 `gorm:"primaryKey;table:resource_tags"`
	TagID      int64 `gorm:"primaryKey"`
}

// ResourceComment 资源评论模型
type ResourceComment struct {
	CommentID  int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"not null"`
	ResourceID int64     `gorm:"not null"`
	Content    string    `gorm:"type:text;not null"`
	ParentID   *int64    `gorm:"default:NULL"`
	Likes      int64     `gorm:"default:0"`
	IsVisible  bool      `gorm:"default:true"`
	Status     string    `gorm:"type:enum('normal','deleted_by_user','deleted_by_admin');default:'normal'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	// 关联用户信息
	User User `gorm:"foreignKey:UserID;references:UserID"`
}

// ResourceRating 资源评分模型
type ResourceRating struct {
	RatingID       int64     `gorm:"primaryKey;autoIncrement"`
	UserID         int64     `gorm:"not null"`
	ResourceID     int64     `gorm:"not null"`
	Recommendation float64   `gorm:"type:decimal(2,1);not null"`
	IsVisible      bool      `gorm:"default:true"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	// 关联用户信息
	User User `gorm:"foreignKey:UserID;references:UserID"`
	// 关联资源信息
	Resource Resource `gorm:"foreignKey:ResourceID;references:ResourceID"`
}

// Review 审核模型
type Review struct {
	ReviewID   int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"not null"`
	TargetID   int64     `gorm:"not null"`
	TargetType string    `gorm:"size:50;not null"`
	Reason     string    `gorm:"type:text;not null"`
	Status     string    `gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
