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

// 保留你的Course相关结构体
type Course struct {
	CourseID    int64     `json:"course_id" db:"course_id"`
	CourseName  string    `json:"course_name" db:"course_name"`
	TeacherID   int64     `json:"teacher_id" db:"teacher_id"`
	Credit      float64   `json:"credit" db:"credit"`
	MajorID     int64     `json:"major_id" db:"major_id"`
	Grade       string    `json:"grade" db:"grade"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (c Course) ToCourseModule() *module.Course {
	course := &module.Course{
		CourseId:   c.CourseID,
		CourseName: c.CourseName,
		TeacherId:  c.TeacherID,
		Credit:     c.Credit,
		MajorId:    c.MajorID,
		Grade:      c.Grade,
		CreatedAt:  c.CreatedAt.Unix(),
		UpdatedAt:  c.UpdatedAt.Unix(),
	}

	// 确保description有合理值
	if c.Description != nil {
		course.Description = *c.Description
	} else {
		course.Description = "暂无描述" // 提供默认值
	}
	return course
}

// CourseRating 课程评分
type CourseRating struct {
	RatingID       int64     `json:"rating_id" db:"rating_id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	CourseID       int64     `json:"course_id" db:"course_id"`
	Recommendation int64     `json:"recommendation" db:"recommendation"`
	Difficulty     string    `json:"difficulty" db:"difficulty"`
	Workload       int64     `json:"workload" db:"workload"`
	Usefulness     int64     `json:"usefulness" db:"usefulness"`
	IsVisible      bool      `json:"is_visible" db:"is_visible"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func (r CourseRating) ToCourseRatingModule() *module.CourseRating {
	return &module.CourseRating{
		RatingId:       r.RatingID,
		UserId:         r.UserID,
		CourseId:       r.CourseID,
		Recommendation: r.Recommendation,
		Difficulty:     r.Difficulty,
		Workload:       r.Workload,
		Usefulness:     r.Usefulness,
		IsVisible:      r.IsVisible,
	}
}

// CourseComment 课程评论
type CourseComment struct {
	CommentID int64     `json:"comment_id" db:"comment_id"`
	CourseID  int64     `json:"course_id" db:"course_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	ParentID  int64     `json:"parent_id" db:"parent_id"`
	IsVisible bool      `json:"is_visible" db:"is_visible"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (c CourseComment) ToCourseCommentModule() *module.CourseComment {
	return &module.CourseComment{
		CommentId: c.CommentID,
		UserId:    c.UserID,
		CourseId:  c.CourseID,
		Content:   c.Content,
		ParentId:  c.ParentID,
		Likes:     0, // 必须：Thrift required字段
		IsVisible: c.IsVisible,
		Status:    0,                  // 必须：Thrift required字段
		CreatedAt: c.CreatedAt.Unix(), // 必须：时间戳转换
	}
}

// 添加上游的Resource相关结构体
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

// 添加上游的转换方法
// ToResourceModule 将db.Resource转换为model.Resource
func (r Resource) ToResourceModule() *module.Resource {
	var tags []*module.ResourceTag
	for _, t := range r.Tags {
		tags = append(tags, t.ToResourceTagModule())
	}

	return &module.Resource{
		ResourceId:    r.ResourceID,
		Title:         r.Title,
		Description:   &r.Description,
		FilePath:      r.FilePath,
		FileType:      r.FileType,
		FileSize:      r.FileSize,
		UploaderId:    r.UploaderID,
		CourseId:      r.CourseID,
		DownloadCount: r.DownloadCount,
		AverageRating: r.AverageRating,
		RatingCount:   r.RatingCount,
		Status:        r.Status,
		CreatedAt:     r.CreatedAt.Unix(),
		Tags:          tags,
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
