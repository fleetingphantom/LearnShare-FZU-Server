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
		Credit:     c.Credit, // 根据 Thrift 定义，Credit 应该是 float64
		MajorId:    c.MajorID,
		Grade:      c.Grade,
	}

	if c.Description != nil {
		course.Description = *c.Description
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
		CourseId:  c.CourseID,
		UserId:    c.UserID,
		Content:   c.Content,
		ParentId:  c.ParentID,
		IsVisible: c.IsVisible,
	}
}

type Resource struct {
	ResourceID    int64     `json:"resource_id" db:"resource_id"`
	CourseID      int64     `json:"course_id" db:"course_id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	Title         string    `json:"title" db:"title"`
	Description   *string   `json:"description,omitempty" db:"description"`
	Type          string    `json:"type" db:"type"`
	FileURL       string    `json:"file_url" db:"file_url"`
	FileSize      int64     `json:"file_size" db:"file_size"`
	Status        string    `json:"status" db:"status"`
	DownloadCount int64     `json:"download_count" db:"download_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

func (r Resource) ToResourceModule() *module.Resource {
	resource := &module.Resource{
		ResourceId:    r.ResourceID,
		CourseId:      r.CourseID,
		UploaderId:    r.UserID,
		Title:         r.Title,
		FilePath:      r.FileURL,
		FileType:      r.Type,
		FileSize:      r.FileSize,
		DownloadCount: r.DownloadCount,
		AverageRating: 0.0,
		RatingCount:   0,
		Status:        int32(0),
		CreatedAt:     r.CreatedAt.Unix(),
	}
	
	resource.Description = r.Description

	return resource
}
