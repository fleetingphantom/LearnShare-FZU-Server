package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/resource"
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupResourceServiceTestDB 初始化资源服务模块测试数据库
func setupResourceServiceTestDB(t *testing.T) func() {
	t.Helper()
	sqliteDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("初始化SQLite失败: %v", err)
	}

	// 创建资源表
	createResourceTableSQL := `
CREATE TABLE IF NOT EXISTS resources (
    resource_id INTEGER PRIMARY KEY AUTOINCREMENT,
    resource_name TEXT NOT NULL,
    description TEXT,
    resource_url TEXT NOT NULL,
    type TEXT NOT NULL,
    size INTEGER NOT NULL,
    uploader_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    download_count INTEGER DEFAULT 0,
    average_rating REAL DEFAULT 0.0,
    rating_count INTEGER DEFAULT 0,
    status TEXT DEFAULT 'pending_review',
    created_at DATETIME
);
`

	// 创建标签表
	createTagTableSQL := `
CREATE TABLE IF NOT EXISTS tags (
    tag_id INTEGER PRIMARY KEY AUTOINCREMENT,
    tag_name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME
);
`

	// 创建资源标签映射表
	createResourceTagMappingSQL := `
CREATE TABLE IF NOT EXISTS resource_tags (
    resource_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (resource_id, tag_id)
);
`

	// 创建资源评论表
	createCommentTableSQL := `
CREATE TABLE IF NOT EXISTS resource_comments (
    comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    parent_id INTEGER,
    likes INTEGER DEFAULT 0,
    is_visible INTEGER DEFAULT 1,
    status TEXT DEFAULT 'normal',
    created_at DATETIME
);
`

	// 创建资源评分表
	createRatingTableSQL := `
CREATE TABLE IF NOT EXISTS resource_ratings (
    rating_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,
    recommendation REAL NOT NULL,
    is_visible INTEGER DEFAULT 1,
    created_at DATETIME
);
`

	// 创建用户表
	createUserTableSQL := `
CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    password_hash TEXT,
    email TEXT UNIQUE,
    college_id INTEGER,
    major_id INTEGER,
    avatar_url TEXT,
    reputation_score INTEGER DEFAULT 0,
    role_id INTEGER,
    status TEXT,
    created_at DATETIME,
    updated_at DATETIME
);
`

	// 创建审核表
	createReviewTableSQL := `
CREATE TABLE IF NOT EXISTS reviews (
    review_id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_id INTEGER NOT NULL,
    target_type TEXT NOT NULL,
    reason TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    priority INTEGER DEFAULT 3,
    reviewer_id INTEGER,
    reviewed_at DATETIME,
    created_at DATETIME
);
`

	tables := []string{
		createResourceTableSQL,
		createTagTableSQL,
		createResourceTagMappingSQL,
		createCommentTableSQL,
		createRatingTableSQL,
		createUserTableSQL,
		createReviewTableSQL,
	}

	for _, sql := range tables {
		if err := sqliteDB.Exec(sql).Error; err != nil {
			t.Fatalf("创建测试数据表失败: %v", err)
		}
	}

	db.DB = sqliteDB

	return func() {
		sqlDB, err := db.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

// seedResourceForService 为service层测试插入资源数据
func seedResourceForService(t *testing.T, name, description string, courseID int64) *db.Resource {
	t.Helper()
	now := time.Now()
	resourcedata := &db.Resource{
		ResourceName:  name,
		Description:   description,
		FilePath:      "/files/test.pdf",
		FileType:      "pdf",
		FileSize:      1024,
		UploaderID:    1,
		CourseID:      courseID,
		DownloadCount: 0,
		AverageRating: 0.0,
		RatingCount:   0,
		Status:        "normal",
		CreatedAt:     now,
	}

	if err := db.DB.WithContext(context.Background()).Table(constants.ResourceTableName).Create(resourcedata).Error; err != nil {
		t.Fatalf("插入测试资源失败: %v", err)
	}
	return resourcedata
}

// seedUserForService 为service层测试插入用户数据
func seedUserForService(t *testing.T, username, email string) *db.User {
	t.Helper()
	user := &db.User{
		Username:        username,
		PasswordHash:    "hash",
		Email:           email,
		ReputationScore: 0,
		RoleID:          2,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := db.DB.WithContext(context.Background()).Table(constants.UserTableName).Create(user).Error; err != nil {
		t.Fatalf("插入测试用户失败: %v", err)
	}
	return user
}

// buildTestRequestContext 构建测试请求上下文
func buildTestRequestContext(userID int64) *app.RequestContext {
	ctx := app.NewContext(0)
	ctx.Set(constants.ContextUid, userID)
	return ctx
}

func TestResourceServiceSearchResources(t *testing.T) {
	cleanup := setupResourceServiceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	c := buildTestRequestContext(1)

	// 准备测试数据
	seedResourceForService(t, "数据结构教程", "经典数据结构教材", 1)
	seedResourceForService(t, "算法导论", "算法入门书籍", 1)
	seedResourceForService(t, "操作系统", "计算机操作系统", 2)

	svc := NewResourceService(ctx, c)

	t.Run("无过滤条件搜索", func(t *testing.T) {
		req := &resource.SearchResourceReq{
			PageNum:  1,
			PageSize: 10,
		}
		resources, total, err := svc.SearchResources(req)
		if err != nil {
			t.Fatalf("搜索资源失败: %v", err)
		}
		if total != 3 {
			t.Fatalf("预期返回3条资源，实际返回 %d", total)
		}
		if len(resources) != 3 {
			t.Fatalf("预期返回3条资源数据，实际返回 %d", len(resources))
		}
	})

	t.Run("关键词搜索", func(t *testing.T) {
		keyword := "数据结构"
		req := &resource.SearchResourceReq{
			Keyword:  &keyword,
			PageNum:  1,
			PageSize: 10,
		}
		resources, total, err := svc.SearchResources(req)
		if err != nil {
			t.Fatalf("关键词搜索失败: %v", err)
		}
		if total != 1 {
			t.Fatalf("预期返回1条资源，实际返回 %d", total)
		}
		if resources[0].Title != "数据结构教程" {
			t.Fatalf("预期返回'数据结构教程'，实际返回 %s", resources[0].Title)
		}
	})

	t.Run("按课程ID过滤", func(t *testing.T) {
		courseID := int64(1)
		req := &resource.SearchResourceReq{
			CourseID: &courseID,
			PageNum:  1,
			PageSize: 10,
		}
		_, total, err := svc.SearchResources(req)
		if err != nil {
			t.Fatalf("按课程ID过滤失败: %v", err)
		}
		if total != 2 {
			t.Fatalf("预期返回2条资源，实际返回 %d", total)
		}
	})

	t.Run("关键词过长验证", func(t *testing.T) {
		longKeyword := string(make([]byte, 101))
		req := &resource.SearchResourceReq{
			Keyword:  &longKeyword,
			PageNum:  1,
			PageSize: 10,
		}
		_, _, err := svc.SearchResources(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
		if err != errno.ValidationKeywordTooLongError {
			t.Fatalf("预期返回关键词过长错误，实际错误为 %v", err)
		}
	})

	t.Run("分页参数自动修正", func(t *testing.T) {
		req := &resource.SearchResourceReq{
			PageNum:  0,
			PageSize: 0,
		}
		resources, _, err := svc.SearchResources(req)
		if err != nil {
			t.Fatalf("搜索资源失败: %v", err)
		}
		if len(resources) > 20 {
			t.Fatalf("预期返回不超过20条数据，实际返回 %d", len(resources))
		}
	})
}

func TestResourceServiceGetResource(t *testing.T) {
	cleanup := setupResourceServiceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	c := buildTestRequestContext(1)

	testResource := seedResourceForService(t, "测试资源", "测试描述", 1)

	svc := NewResourceService(ctx, c)

	t.Run("成功获取资源", func(t *testing.T) {
		req := &resource.GetResourceReq{
			ResourceID: testResource.ResourceID,
		}
		result, err := svc.GetResource(req)
		if err != nil {
			t.Fatalf("获取资源失败: %v", err)
		}
		if result.ResourceId != testResource.ResourceID {
			t.Fatalf("预期资源ID为 %d，实际为 %d", testResource.ResourceID, result.ResourceId)
		}
		if result.Title != "测试资源" {
			t.Fatalf("预期资源名称为'测试资源'，实际为 %s", result.Title)
		}
	})

	t.Run("资源ID无效", func(t *testing.T) {
		req := &resource.GetResourceReq{
			ResourceID: 0,
		}
		_, err := svc.GetResource(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
	})

	t.Run("资源不存在", func(t *testing.T) {
		req := &resource.GetResourceReq{
			ResourceID: 99999,
		}
		_, err := svc.GetResource(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
	})
}

func TestResourceServiceGetResourceComments(t *testing.T) {
	cleanup := setupResourceServiceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	c := buildTestRequestContext(1)

	testResource := seedResourceForService(t, "测试资源", "测试描述", 1)
	user := seedUserForService(t, "testuser", "test@example.com")

	// 插入评论
	comment1 := &db.ResourceComment{
		UserID:     user.UserID,
		ResourceID: testResource.ResourceID,
		Content:    "好评",
		Likes:      10,
		IsVisible:  true,
		Status:     "normal",
		CreatedAt:  time.Now().Add(-2 * time.Hour),
	}
	comment2 := &db.ResourceComment{
		UserID:     user.UserID,
		ResourceID: testResource.ResourceID,
		Content:    "很好",
		Likes:      5,
		IsVisible:  true,
		Status:     "normal",
		CreatedAt:  time.Now().Add(-1 * time.Hour),
	}

	db.DB.WithContext(ctx).Table(constants.ResourceCommentTableName).Create(comment1)
	db.DB.WithContext(ctx).Table(constants.ResourceCommentTableName).Create(comment2)

	svc := NewResourceService(ctx, c)

	t.Run("获取评论列表", func(t *testing.T) {
		req := &resource.GetResourceCommentsReq{
			ResourceID: testResource.ResourceID,
			PageNum:    1,
			PageSize:   10,
		}
		comments, total, err := svc.GetResourceComments(req)
		if err != nil {
			t.Fatalf("获取评论列表失败: %v", err)
		}
		if total != 2 {
			t.Fatalf("预期返回2条评论，实际返回 %d", total)
		}
		if len(comments) != 2 {
			t.Fatalf("预期返回2条评论数据，实际返回 %d", len(comments))
		}
	})

	t.Run("资源ID无效", func(t *testing.T) {
		req := &resource.GetResourceCommentsReq{
			ResourceID: 0,
			PageNum:    1,
			PageSize:   10,
		}
		_, _, err := svc.GetResourceComments(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
	})

	t.Run("分页参数自动修正", func(t *testing.T) {
		req := &resource.GetResourceCommentsReq{
			ResourceID: testResource.ResourceID,
			PageNum:    0,
			PageSize:   0,
		}
		comments, _, err := svc.GetResourceComments(req)
		if err != nil {
			t.Fatalf("获取评论列表失败: %v", err)
		}
		if len(comments) > 20 {
			t.Fatalf("预期返回不超过20条数据，实际返回 %d", len(comments))
		}
	})
}

func TestResourceServiceSubmitResourceRating(t *testing.T) {
	cleanup := setupResourceServiceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	user := seedUserForService(t, "testuser", "test@example.com")
	c := buildTestRequestContext(user.UserID)

	testResource := seedResourceForService(t, "测试资源", "测试描述", 1)

	svc := NewResourceService(ctx, c)

	t.Run("提交评分成功", func(t *testing.T) {
		req := &resource.SubmitResourceRatingReq{
			ResourceID: testResource.ResourceID,
			Rating:     4.5,
		}
		rating, err := svc.SubmitResourceRating(req)
		if err != nil {
			t.Fatalf("提交评分失败: %v", err)
		}
		// Recommendation被转换为0-50的范围（乘以10）
		if rating.Recommendation != 45.0 {
			t.Fatalf("预期评分为45.0（4.5*10），实际为 %f", rating.Recommendation)
		}
	})

	t.Run("评分范围验证-小于0", func(t *testing.T) {
		req := &resource.SubmitResourceRatingReq{
			ResourceID: testResource.ResourceID,
			Rating:     -1.0,
		}
		_, err := svc.SubmitResourceRating(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
		if err != errno.ValidationRatingRangeInvalidError {
			t.Fatalf("预期返回评分范围错误，实际错误为 %v", err)
		}
	})

	t.Run("评分范围验证-大于5", func(t *testing.T) {
		req := &resource.SubmitResourceRatingReq{
			ResourceID: testResource.ResourceID,
			Rating:     5.5,
		}
		_, err := svc.SubmitResourceRating(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
		if err != errno.ValidationRatingRangeInvalidError {
			t.Fatalf("预期返回评分范围错误，实际错误为 %v", err)
		}
	})
}

func TestResourceServiceSubmitResourceComment(t *testing.T) {
	cleanup := setupResourceServiceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	user := seedUserForService(t, "testuser", "test@example.com")
	c := buildTestRequestContext(user.UserID)

	testResource := seedResourceForService(t, "测试资源", "测试描述", 1)

	svc := NewResourceService(ctx, c)

	t.Run("提交评论成功", func(t *testing.T) {
		req := &resource.SubmitResourceCommentReq{
			ResourceID: testResource.ResourceID,
			Content:    "这是一条评论",
		}
		comment, err := svc.SubmitResourceComment(req)
		if err != nil {
			t.Fatalf("提交评论失败: %v", err)
		}
		if comment.Content != "这是一条评论" {
			t.Fatalf("预期评论内容为'这是一条评论'，实际为 %s", comment.Content)
		}
	})

	t.Run("评论内容为空", func(t *testing.T) {
		req := &resource.SubmitResourceCommentReq{
			ResourceID: testResource.ResourceID,
			Content:    "",
		}
		_, err := svc.SubmitResourceComment(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
		if err != errno.ResourceInvalidCommentError {
			t.Fatalf("预期返回评论无效错误，实际错误为 %v", err)
		}
	})

	t.Run("评论内容过长", func(t *testing.T) {
		longContent := string(make([]byte, 1001))
		req := &resource.SubmitResourceCommentReq{
			ResourceID: testResource.ResourceID,
			Content:    longContent,
		}
		_, err := svc.SubmitResourceComment(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
		if err != errno.ValidationCommentTooLongError {
			t.Fatalf("预期返回评论过长错误，实际错误为 %v", err)
		}
	})

	t.Run("提交回复评论", func(t *testing.T) {
		parentID := int64(1)
		req := &resource.SubmitResourceCommentReq{
			ResourceID: testResource.ResourceID,
			Content:    "这是一条回复",
			ParentId:   &parentID,
		}
		comment, err := svc.SubmitResourceComment(req)
		if err != nil {
			t.Fatalf("提交回复评论失败: %v", err)
		}
		if comment.ParentId != parentID {
			t.Fatalf("预期父评论ID为 %d，实际为 %d", parentID, comment.ParentId)
		}
	})
}

func TestResourceServiceDeleteResourceRating(t *testing.T) {
	cleanup := setupResourceServiceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	user := seedUserForService(t, "testuser", "test@example.com")
	c := buildTestRequestContext(user.UserID)

	testResource := seedResourceForService(t, "测试资源", "测试描述", 1)

	svc := NewResourceService(ctx, c)

	// 先提交一个评分
	submitReq := &resource.SubmitResourceRatingReq{
		ResourceID: testResource.ResourceID,
		Rating:     4.0,
	}
	rating, err := svc.SubmitResourceRating(submitReq)
	if err != nil {
		t.Fatalf("提交评分失败: %v", err)
	}

	t.Run("删除评分成功", func(t *testing.T) {
		req := &resource.DeleteResourceRatingReq{
			RatingID: rating.RatingId,
		}
		err := svc.DeleteResourceRating(req)
		if err != nil {
			t.Fatalf("删除评分失败: %v", err)
		}
	})

	t.Run("评分ID无效", func(t *testing.T) {
		req := &resource.DeleteResourceRatingReq{
			RatingID: 0,
		}
		err := svc.DeleteResourceRating(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
	})
}

func TestResourceServiceDeleteResourceComment(t *testing.T) {
	cleanup := setupResourceServiceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	user := seedUserForService(t, "testuser", "test@example.com")
	c := buildTestRequestContext(user.UserID)

	testResource := seedResourceForService(t, "测试资源", "测试描述", 1)

	svc := NewResourceService(ctx, c)

	// 先提交一个评论
	submitReq := &resource.SubmitResourceCommentReq{
		ResourceID: testResource.ResourceID,
		Content:    "测试评论",
	}
	comment, err := svc.SubmitResourceComment(submitReq)
	if err != nil {
		t.Fatalf("提交评论失败: %v", err)
	}

	t.Run("删除评论成功", func(t *testing.T) {
		req := &resource.DeleteResourceCommentReq{
			CommentID: comment.CommentId,
		}
		err := svc.DeleteResourceComment(req)
		if err != nil {
			t.Fatalf("删除评论失败: %v", err)
		}
	})

	t.Run("评论ID无效", func(t *testing.T) {
		req := &resource.DeleteResourceCommentReq{
			CommentID: 0,
		}
		err := svc.DeleteResourceComment(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
	})
}

func TestResourceServiceReportResource(t *testing.T) {
	cleanup := setupResourceServiceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	user := seedUserForService(t, "testuser", "test@example.com")
	c := buildTestRequestContext(user.UserID)

	testResource := seedResourceForService(t, "测试资源", "测试描述", 1)

	svc := NewResourceService(ctx, c)

	t.Run("举报资源成功", func(t *testing.T) {
		req := &resource.ReportResourceReq{
			ResourceID: testResource.ResourceID,
			Reason:     "违规内容",
		}
		err := svc.ReportResource(req)
		if err != nil {
			t.Fatalf("举报资源失败: %v", err)
		}
	})

	t.Run("资源ID无效", func(t *testing.T) {
		req := &resource.ReportResourceReq{
			ResourceID: 0,
			Reason:     "违规内容",
		}
		err := svc.ReportResource(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
	})

	t.Run("举报原因为空", func(t *testing.T) {
		req := &resource.ReportResourceReq{
			ResourceID: testResource.ResourceID,
			Reason:     "",
		}
		err := svc.ReportResource(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
		if err != errno.ResourceReportInvalidReasonError {
			t.Fatalf("预期返回举报原因无效错误，实际错误为 %v", err)
		}
	})

	t.Run("举报原因过长", func(t *testing.T) {
		longReason := string(make([]byte, 501))
		req := &resource.ReportResourceReq{
			ResourceID: testResource.ResourceID,
			Reason:     longReason,
		}
		err := svc.ReportResource(req)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
		if err != errno.ValidationReportReasonTooLongError {
			t.Fatalf("预期返回举报原因过长错误，实际错误为 %v", err)
		}
	})
}
