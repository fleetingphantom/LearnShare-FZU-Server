package db

import (
	"LearnShare/pkg/constants"
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupResourceTestDB 初始化资源模块测试数据库
func setupResourceTestDB(t *testing.T) func() {
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

	DB = sqliteDB

	return func() {
		sqlDB, err := DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

// seedResource 插入测试资源数据
func seedResource(t *testing.T, name, description string, courseID int64) *Resource {
	t.Helper()
	now := time.Now()
	resource := &Resource{
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

	if err := DB.WithContext(context.Background()).Table(constants.ResourceTableName).Create(resource).Error; err != nil {
		t.Fatalf("插入测试资源失败: %v", err)
	}
	return resource
}

// seedTag 插入测试标签数据
func seedTag(t *testing.T, tagName string) *ResourceTag {
	t.Helper()
	tag := &ResourceTag{
		TagName: tagName,
	}

	if err := DB.WithContext(context.Background()).Table("tags").Create(tag).Error; err != nil {
		t.Fatalf("插入测试标签失败: %v", err)
	}
	return tag
}

// seedUser 插入测试用户数据
func seedUser(t *testing.T, username, email string) *User {
	t.Helper()
	user := &User{
		Username:        username,
		PasswordHash:    "hash",
		Email:           email,
		ReputationScore: 0,
		RoleID:          2,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := DB.WithContext(context.Background()).Table(constants.UserTableName).Create(user).Error; err != nil {
		t.Fatalf("插入测试用户失败: %v", err)
	}
	return user
}

// linkResourceTag 关联资源和标签
func linkResourceTag(t *testing.T, resourceID, tagID int64) {
	t.Helper()
	sql := `INSERT INTO resource_tags (resource_id, tag_id) VALUES (?, ?)`
	if err := DB.WithContext(context.Background()).Exec(sql, resourceID, tagID).Error; err != nil {
		t.Fatalf("关联资源标签失败: %v", err)
	}
}

func TestSearchResources(t *testing.T) {
	cleanup := setupResourceTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// 准备测试数据
	resource1 := seedResource(t, "数据结构教程", "经典数据结构教材", 1)
	_ = seedResource(t, "算法导论", "算法入门书籍", 1)
	seedResource(t, "操作系统", "计算机操作系统", 2)

	tag := seedTag(t, "教材")
	linkResourceTag(t, resource1.ResourceID, tag.TagID)

	t.Run("无过滤条件搜索", func(t *testing.T) {
		resources, total, err := SearchResources(ctx, nil, nil, nil, nil, 1, 10)
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
		resources, total, err := SearchResources(ctx, &keyword, nil, nil, nil, 1, 10)
		if err != nil {
			t.Fatalf("关键词搜索失败: %v", err)
		}
		if total != 1 {
			t.Fatalf("预期返回1条资源，实际返回 %d", total)
		}
		if resources[0].ResourceName != "数据结构教程" {
			t.Fatalf("预期返回'数据结构教程'，实际返回 %s", resources[0].ResourceName)
		}
	})

	t.Run("按课程ID过滤", func(t *testing.T) {
		courseID := int64(1)
		_, total, err := SearchResources(ctx, nil, nil, &courseID, nil, 1, 10)
		if err != nil {
			t.Fatalf("按课程ID过滤失败: %v", err)
		}
		if total != 2 {
			t.Fatalf("预期返回2条资源，实际返回 %d", total)
		}
	})

	t.Run("按标签ID过滤", func(t *testing.T) {
		tagID := tag.TagID
		resources, total, err := SearchResources(ctx, nil, &tagID, nil, nil, 1, 10)
		if err != nil {
			t.Fatalf("按标签ID过滤失败: %v", err)
		}
		if total != 1 {
			t.Fatalf("预期返回1条资源，实际返回 %d", total)
		}
		if resources[0].ResourceID != resource1.ResourceID {
			t.Fatalf("预期返回resource1，实际返回 %d", resources[0].ResourceID)
		}
	})

	t.Run("分页测试", func(t *testing.T) {
		_, total, err := SearchResources(ctx, nil, nil, nil, nil, 1, 2)
		if err != nil {
			t.Fatalf("分页查询失败: %v", err)
		}
		if total != 3 {
			t.Fatalf("预期总数3，实际返回 %d", total)
		}
	})
}

func TestGetResourceByID(t *testing.T) {
	cleanup := setupResourceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	resource := seedResource(t, "测试资源", "测试描述", 1)

	t.Run("成功获取资源", func(t *testing.T) {
		result, err := GetResourceByID(ctx, resource.ResourceID)
		if err != nil {
			t.Fatalf("获取资源失败: %v", err)
		}
		if result.ResourceID != resource.ResourceID {
			t.Fatalf("预期资源ID为 %d，实际为 %d", resource.ResourceID, result.ResourceID)
		}
		if result.ResourceName != "测试资源" {
			t.Fatalf("预期资源名称为'测试资源'，实际为 %s", result.ResourceName)
		}
	})

	t.Run("资源不存在", func(t *testing.T) {
		_, err := GetResourceByID(ctx, 99999)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
	})
}

func TestGetResourceComments(t *testing.T) {
	cleanup := setupResourceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	resource := seedResource(t, "测试资源", "测试描述", 1)
	user := seedUser(t, "testuser", "test@example.com")

	// 插入评论
	comment1 := &ResourceComment{
		UserID:     user.UserID,
		ResourceID: resource.ResourceID,
		Content:    "好评",
		Likes:      10,
		IsVisible:  true,
		Status:     "normal",
		CreatedAt:  time.Now().Add(-2 * time.Hour),
	}
	comment2 := &ResourceComment{
		UserID:     user.UserID,
		ResourceID: resource.ResourceID,
		Content:    "很好",
		Likes:      5,
		IsVisible:  true,
		Status:     "normal",
		CreatedAt:  time.Now().Add(-1 * time.Hour),
	}

	DB.WithContext(ctx).Table(constants.ResourceCommentTableName).Create(comment1)
	DB.WithContext(ctx).Table(constants.ResourceCommentTableName).Create(comment2)

	t.Run("获取评论列表", func(t *testing.T) {
		comments, total, err := GetResourceComments(ctx, resource.ResourceID, nil, 1, 10)
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

	t.Run("按最新排序", func(t *testing.T) {
		sortBy := "latest"
		comments, _, err := GetResourceComments(ctx, resource.ResourceID, &sortBy, 1, 10)
		if err != nil {
			t.Fatalf("按最新排序失败: %v", err)
		}
		if comments[0].Content != "很好" {
			t.Fatalf("预期第一条评论为'很好'，实际为 %s", comments[0].Content)
		}
	})

	t.Run("按热度排序", func(t *testing.T) {
		sortBy := "hottest"
		comments, _, err := GetResourceComments(ctx, resource.ResourceID, &sortBy, 1, 10)
		if err != nil {
			t.Fatalf("按热度排序失败: %v", err)
		}
		if comments[0].Content != "好评" {
			t.Fatalf("预期第一条评论为'好评'，实际为 %s", comments[0].Content)
		}
	})
}

func TestSubmitResourceRating(t *testing.T) {
	cleanup := setupResourceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	resource := seedResource(t, "测试资源", "测试描述", 1)

	t.Run("提交新评分", func(t *testing.T) {
		rating, err := SubmitResourceRating(ctx, 1, resource.ResourceID, 4.5)
		if err != nil {
			t.Fatalf("提交评分失败: %v", err)
		}
		if rating.Recommendation != 4.5 {
			t.Fatalf("预期评分为4.5，实际为 %f", rating.Recommendation)
		}

		// 验证资源的评分信息已更新
		updatedResource, err := GetResourceByID(ctx, resource.ResourceID)
		if err != nil {
			t.Fatalf("获取资源失败: %v", err)
		}
		if updatedResource.AverageRating != 4.5 {
			t.Fatalf("预期平均评分为4.5，实际为 %f", updatedResource.AverageRating)
		}
		if updatedResource.RatingCount != 1 {
			t.Fatalf("预期评分数量为1，实际为 %d", updatedResource.RatingCount)
		}
	})

	t.Run("更新已有评分", func(t *testing.T) {
		// 再次提交评分
		rating, err := SubmitResourceRating(ctx, 1, resource.ResourceID, 5.0)
		if err != nil {
			t.Fatalf("更新评分失败: %v", err)
		}
		if rating.Recommendation != 5.0 {
			t.Fatalf("预期评分为5.0，实际为 %f", rating.Recommendation)
		}

		// 验证资源的评分信息已更新
		updatedResource, err := GetResourceByID(ctx, resource.ResourceID)
		if err != nil {
			t.Fatalf("获取资源失败: %v", err)
		}
		if updatedResource.AverageRating != 5.0 {
			t.Fatalf("预期平均评分为5.0，实际为 %f", updatedResource.AverageRating)
		}
		if updatedResource.RatingCount != 1 {
			t.Fatalf("预期评分数量仍为1，实际为 %d", updatedResource.RatingCount)
		}
	})
}

func TestSubmitResourceComment(t *testing.T) {
	cleanup := setupResourceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	resource := seedResource(t, "测试资源", "测试描述", 1)
	user := seedUser(t, "testuser", "test@example.com")

	t.Run("提交评论", func(t *testing.T) {
		comment, err := SubmitResourceComment(ctx, user.UserID, resource.ResourceID, "这是一条评论", nil)
		if err != nil {
			t.Fatalf("提交评论失败: %v", err)
		}
		if comment.Content != "这是一条评论" {
			t.Fatalf("预期评论内容为'这是一条评论'，实际为 %s", comment.Content)
		}
		if comment.Status != "normal" {
			t.Fatalf("预期评论状态为'normal'，实际为 %s", comment.Status)
		}
	})

	t.Run("提交回复评论", func(t *testing.T) {
		parentID := int64(1)
		comment, err := SubmitResourceComment(ctx, user.UserID, resource.ResourceID, "这是一条回复", &parentID)
		if err != nil {
			t.Fatalf("提交回复评论失败: %v", err)
		}
		if comment.ParentID == nil || *comment.ParentID != parentID {
			t.Fatalf("预期父评论ID为 %d，实际为 %v", parentID, comment.ParentID)
		}
	})
}

func TestDeleteResourceRating(t *testing.T) {
	cleanup := setupResourceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	resource := seedResource(t, "测试资源", "测试描述", 1)

	// 先提交一个评分
	rating, err := SubmitResourceRating(ctx, 1, resource.ResourceID, 4.0)
	if err != nil {
		t.Fatalf("提交评分失败: %v", err)
	}

	t.Run("删除评分成功", func(t *testing.T) {
		err := DeleteResourceRating(ctx, rating.RatingID, 1)
		if err != nil {
			t.Fatalf("删除评分失败: %v", err)
		}

		// 验证资源的评分信息已更新
		updatedResource, err := GetResourceByID(ctx, resource.ResourceID)
		if err != nil {
			t.Fatalf("获取资源失败: %v", err)
		}
		if updatedResource.RatingCount != 0 {
			t.Fatalf("预期评分数量为0，实际为 %d", updatedResource.RatingCount)
		}
	})

	t.Run("删除不存在的评分", func(t *testing.T) {
		err := DeleteResourceRating(ctx, 99999, 1)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
	})
}

func TestDeleteResourceComment(t *testing.T) {
	cleanup := setupResourceTestDB(t)
	defer cleanup()

	ctx := context.Background()
	resource := seedResource(t, "测试资源", "测试描述", 1)
	user := seedUser(t, "testuser", "test@example.com")

	// 先提交一个评论
	comment, err := SubmitResourceComment(ctx, user.UserID, resource.ResourceID, "测试评论", nil)
	if err != nil {
		t.Fatalf("提交评论失败: %v", err)
	}

	t.Run("删除评论成功", func(t *testing.T) {
		err := DeleteResourceComment(ctx, comment.CommentID, user.UserID)
		if err != nil {
			t.Fatalf("删除评论失败: %v", err)
		}

		// 验证评论已删除
		comments, total, err := GetResourceComments(ctx, resource.ResourceID, nil, 1, 10)
		if err != nil {
			t.Fatalf("获取评论列表失败: %v", err)
		}
		if total != 0 {
			t.Fatalf("预期评论数量为0，实际为 %d", total)
		}
		if len(comments) != 0 {
			t.Fatalf("预期返回0条评论，实际返回 %d", len(comments))
		}
	})

	t.Run("删除不存在的评论", func(t *testing.T) {
		err := DeleteResourceComment(ctx, 99999, user.UserID)
		if err == nil {
			t.Fatalf("预期返回错误，实际成功")
		}
	})
}

func TestCreateReview(t *testing.T) {
	cleanup := setupResourceTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("创建举报成功", func(t *testing.T) {
		err := CreateReview(ctx, 1, 100, "resource", "违规内容")
		if err != nil {
			t.Fatalf("创建举报失败: %v", err)
		}

		// 验证举报记录已创建
		var review Review
		err = DB.WithContext(ctx).Table(constants.ReviewTableName).Where("target_id = ? AND target_type = ?", 100, "resource").First(&review).Error
		if err != nil {
			t.Fatalf("查询举报记录失败: %v", err)
		}
		if review.Reason != "违规内容" {
			t.Fatalf("预期举报原因为'违规内容'，实际为 %s", review.Reason)
		}
	})
}
