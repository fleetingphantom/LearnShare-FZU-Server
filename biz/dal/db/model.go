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

type Resource struct {
	ResourceID    int64         `gorm:"primaryKey;autoIncrement"`
	Title         string        `gorm:"size:255;not null"`
	Description   string        `gorm:"type:text"`
	FilePath      string        `gorm:"size:255;not null"`
	FileType      string        `gorm:"size:50;not null"`
	FileSize      int64         `gorm:"not null"`
	UploaderID    int64         `gorm:"not null"`
	CourseID      int64         `gorm:"not null"`
	DownloadCount int64         `gorm:"default:0"`
	AverageRating float64       `gorm:"default:0.0"`
	RatingCount   int64         `gorm:"default:0"`
	Status        int32         `gorm:"not null;default:0"`
	CreatedAt     time.Time     `gorm:"autoCreateTime"`
	Tags          []ResourceTag `gorm:"many2many:resource_tag_mappings;"`
}

type ResourceTag struct {
	TagID   int64  `gorm:"primaryKey;autoIncrement"`
	TagName string `gorm:"size:50;unique;not null"`
}

type ResourceTagMapping struct {
	ResourceID int64 `gorm:"primaryKey"`
	TagID      int64 `gorm:"primaryKey"`
}
