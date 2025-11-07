package course

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"LearnShare/biz/dal/db"
	courseModel "LearnShare/biz/model/course"
	"LearnShare/pkg/constants"

	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupHandlerTestDB 初始化 handler 测试数据库
func setupHandlerTestDB(t *testing.T) func() {
	t.Helper()
	sqliteDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("初始化SQLite失败: %v", err)
	}

	// 创建课程表
	createCourseTableSQL := `
CREATE TABLE IF NOT EXISTS courses (
    course_id INTEGER PRIMARY KEY AUTOINCREMENT,
    course_name TEXT NOT NULL,
    teacher_id INTEGER NOT NULL,
    credit REAL NOT NULL,
    major_id INTEGER NOT NULL,
    grade TEXT NOT NULL,
    description TEXT,
    created_at DATETIME,
    updated_at DATETIME
);
`
	if err := sqliteDB.Exec(createCourseTableSQL).Error; err != nil {
		t.Fatalf("创建课程表失败: %v", err)
	}

	// 创建课程评分表
	createCourseRatingTableSQL := `
CREATE TABLE IF NOT EXISTS course_ratings (
    rating_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    recommendation INTEGER NOT NULL,
    difficulty TEXT NOT NULL,
    workload INTEGER NOT NULL,
    usefulness INTEGER NOT NULL,
    is_visible BOOLEAN DEFAULT 1,
    created_at DATETIME,
    updated_at DATETIME
);
`
	if err := sqliteDB.Exec(createCourseRatingTableSQL).Error; err != nil {
		t.Fatalf("创建课程评分表失败: %v", err)
	}

	// 创建课程评论表
	createCourseCommentTableSQL := `
CREATE TABLE IF NOT EXISTS course_comments (
    comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    course_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    parent_id INTEGER DEFAULT 0,
    is_visible BOOLEAN DEFAULT 1,
    created_at DATETIME,
    updated_at DATETIME
);
`
	if err := sqliteDB.Exec(createCourseCommentTableSQL).Error; err != nil {
		t.Fatalf("创建课程评论表失败: %v", err)
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
	if err := sqliteDB.Exec(createResourceTableSQL).Error; err != nil {
		t.Fatalf("创建资源表失败: %v", err)
	}

	db.DB = sqliteDB

	return func() {
		sqlDB, err := db.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

// seedHandlerCourse 插入测试课程
func seedHandlerCourse(t *testing.T, courseName string, teacherID, majorID int64, grade string) *db.Course {
	t.Helper()
	now := time.Now()
	description := "测试课程描述"
	courseRecord := &db.Course{
		CourseName:  courseName,
		TeacherID:   teacherID,
		Credit:      3.0,
		MajorID:     majorID,
		Grade:       grade,
		Description: &description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := db.DB.WithContext(context.Background()).Table(constants.CourseTableName).Create(courseRecord).Error; err != nil {
		t.Fatalf("插入测试课程失败: %v", err)
	}
	return courseRecord
}

// TestSearch 测试搜索课程 handler
func TestSearch(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	// 插入测试数据
	seedHandlerCourse(t, "高等数学", 101, 1, "2024")
	seedHandlerCourse(t, "线性代数", 102, 1, "2024")

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/api/courses/search", Search)

	w := ut.PerformRequest(router, "GET", "/api/courses/search?keywords=数学&page_num=1&page_size=10", nil)
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}

	var searchResp courseModel.SearchResp
	if err := json.Unmarshal(resp.Body(), &searchResp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if searchResp.BaseResponse.Code != 0 {
		t.Errorf("期望响应码为 0, 实际为 %d, 消息: %s", searchResp.BaseResponse.Code, searchResp.BaseResponse.Message)
	}

	if len(searchResp.Courses) != 1 {
		t.Errorf("期望课程数量为 1, 实际为 %d", len(searchResp.Courses))
	}
}

// TestSearchNoResults 测试搜索无结果
func TestSearchNoResults(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/api/courses/search", Search)

	w := ut.PerformRequest(router, "GET", "/api/courses/search?keywords=不存在的课程&page_num=1&page_size=10", nil)
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}

	var searchResp courseModel.SearchResp
	if err := json.Unmarshal(resp.Body(), &searchResp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if len(searchResp.Courses) != 0 {
		t.Errorf("期望课程数量为 0, 实际为 %d", len(searchResp.Courses))
	}
}

// TestGetCourseDetail 测试获取课程详情 handler
func TestGetCourseDetail(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	inserted := seedHandlerCourse(t, "计算机网络", 103, 2, "2024")

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/api/courses/:course_id", GetCourseDetail)

	w := ut.PerformRequest(router, "GET", fmt.Sprintf("/api/courses/%d", inserted.CourseID), nil)
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}

	var detailResp courseModel.GetCourseDetailResp
	if err := json.Unmarshal(resp.Body(), &detailResp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if detailResp.BaseResponse.Code != 0 {
		t.Errorf("期望响应码为 0, 实际为 %d", detailResp.BaseResponse.Code)
	}

	if detailResp.Course.CourseName != "计算机网络" {
		t.Errorf("期望课程名为 计算机网络, 实际为 %s", detailResp.Course.CourseName)
	}
}

// TestGetCourseResourceList 测试获取课程资源列表 handler
func TestGetCourseResourceList(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	ctx := context.Background()
	courseID := int64(1)

	// 插入测试课程
	seedHandlerCourse(t, "测试课程", 105, 1, "2024")

	// 插入测试资源
	resources := []*db.Resource{
		{ResourceName: "课件1", FilePath: "/path/1.pdf", FileType: "pdf", FileSize: 1024, UploaderID: 201, CourseID: courseID, Status: "normal", CreatedAt: time.Now()},
		{ResourceName: "课件2", FilePath: "/path/2.docx", FileType: "docx", FileSize: 2048, UploaderID: 202, CourseID: courseID, Status: "normal", CreatedAt: time.Now()},
	}
	for _, r := range resources {
		if err := db.DB.WithContext(ctx).Table(constants.ResourceTableName).Create(r).Error; err != nil {
			t.Fatalf("插入测试资源失败: %v", err)
		}
	}

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/api/courses/:course_id/resources", GetCourseResourceList)

	w := ut.PerformRequest(router, "GET", "/api/courses/1/resources?page_num=1&page_size=10", nil)
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}

	var resourceResp courseModel.GetCourseResourceListResp
	if err := json.Unmarshal(resp.Body(), &resourceResp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resourceResp.BaseResponse.Code != 0 {
		t.Errorf("期望响应码为 0, 实际为 %d", resourceResp.BaseResponse.Code)
	}

	if len(resourceResp.Resources) != 2 {
		t.Errorf("期望资源数量为 2, 实际为 %d", len(resourceResp.Resources))
	}
}

// TestGetCourseComments 测试获取课程评论列表 handler
func TestGetCourseComments(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	ctx := context.Background()
	courseID := int64(1)

	// 插入测试课程
	seedHandlerCourse(t, "测试课程", 106, 1, "2024")

	// 插入测试评论
	comments := []*db.CourseComment{
		{CourseID: courseID, UserID: 301, Content: "很好的课程", ParentID: 0, IsVisible: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CourseID: courseID, UserID: 302, Content: "老师讲得很清楚", ParentID: 0, IsVisible: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, c := range comments {
		if err := db.DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(c).Error; err != nil {
			t.Fatalf("插入测试评论失败: %v", err)
		}
	}

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/api/courses/:course_id/comments", GetCourseComments)

	w := ut.PerformRequest(router, "GET", "/api/courses/1/comments?sort_by=latest&page_num=1&page_size=10", nil)
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}

	var commentResp courseModel.GetCourseCommentsResp
	if err := json.Unmarshal(resp.Body(), &commentResp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if commentResp.BaseResponse.Code != 0 {
		t.Errorf("期望响应码为 0, 实际为 %d", commentResp.BaseResponse.Code)
	}

	if len(commentResp.Comments) != 2 {
		t.Errorf("期望评论数量为 2, 实际为 %d", len(commentResp.Comments))
	}
}

// TestSubmitCourseRating 测试提交课程评分 handler
func TestSubmitCourseRating(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.POST("/api/course_ratings/:course_id", SubmitCourseRating)

	requestBody := map[string]interface{}{
		"course_id": 1,
		"rating":    5,
	}
	body, _ := json.Marshal(requestBody)

	w := ut.PerformRequest(router, "POST", "/api/course_ratings/1", &ut.Body{Body: bytes.NewBuffer(body), Len: len(body)},
		ut.Header{"Content-Type", "application/json"},
		ut.Header{Key: constants.ContextUid, Value: "401"})
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}
}

// TestSubmitCourseComment 测试提交课程评论 handler
func TestSubmitCourseComment(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.POST("/api/courses/:course_id/comments", SubmitCourseComment)

	requestBody := map[string]interface{}{
		"course_id":  1,
		"contents":   "这是一个测试评论",
		"parent_id":  0,
		"is_visible": true,
	}
	body, _ := json.Marshal(requestBody)

	w := ut.PerformRequest(router, "POST", "/api/courses/1/comments", &ut.Body{Body: bytes.NewBuffer(body), Len: len(body)},
		ut.Header{"Content-Type", "application/json"},
		ut.Header{Key: constants.ContextUid, Value: "402"})
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}
}

// TestDeleteCourseComment 测试删除课程评论 handler
func TestDeleteCourseComment(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	ctx := context.Background()
	comment := &db.CourseComment{
		CourseID:  1,
		UserID:    403,
		Content:   "待删除的评论",
		ParentID:  0,
		IsVisible: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(comment).Error; err != nil {
		t.Fatalf("插入测试评论失败: %v", err)
	}

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.DELETE("/api/courses_comments/:comment_id", DeleteCourseComment)

	w := ut.PerformRequest(router, "DELETE", fmt.Sprintf("/api/courses_comments/%d", comment.CommentID), nil)
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}

	var deleteResp courseModel.DeleteCourseCommentResp
	if err := json.Unmarshal(resp.Body(), &deleteResp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if deleteResp.BaseResponse.Code != 0 {
		t.Errorf("期望响应码为 0, 实际为 %d", deleteResp.BaseResponse.Code)
	}
}

// TestDeleteCourseRating 测试删除课程评分 handler
func TestDeleteCourseRating(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rating := &db.CourseRating{
		UserID:         404,
		CourseID:       1,
		Recommendation: 4,
		Difficulty:     "medium",
		Workload:       3,
		Usefulness:     4,
		IsVisible:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := db.DB.WithContext(ctx).Table(constants.CourseRatingTableName).Create(rating).Error; err != nil {
		t.Fatalf("插入测试评分失败: %v", err)
	}

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.DELETE("/api/course_ratings/:rating_id", DeleteCourseRating)

	w := ut.PerformRequest(router, "DELETE", fmt.Sprintf("/api/course_ratings/%d", rating.RatingID), nil)
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}

	var deleteResp courseModel.DeleteCourseRatingResp
	if err := json.Unmarshal(resp.Body(), &deleteResp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if deleteResp.BaseResponse.Code != 0 {
		t.Errorf("期望响应码为 0, 实际为 %d", deleteResp.BaseResponse.Code)
	}
}

// TestGetCourseDetailNotFound 测试获取不存在的课程
func TestGetCourseDetailNotFound(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/api/courses/:course_id", GetCourseDetail)

	w := ut.PerformRequest(router, "GET", "/api/courses/999999", nil)
	resp := w.Result()

	// 应该返回错误或者找不到记录
	if resp.StatusCode() == 200 {
		var detailResp courseModel.GetCourseDetailResp
		if err := json.Unmarshal(resp.Body(), &detailResp); err == nil {
			if detailResp.BaseResponse.Code == 0 {
				t.Error("期望返回错误响应，但返回了成功")
			}
		}
	}
}

// TestGetCourseResourceListWithTypeFilter 测试带类型过滤的资源列表
func TestGetCourseResourceListWithTypeFilter(t *testing.T) {
	cleanup := setupHandlerTestDB(t)
	defer cleanup()

	ctx := context.Background()
	courseID := int64(1)

	// 插入测试课程
	seedHandlerCourse(t, "测试课程", 107, 1, "2024")

	// 插入不同类型的资源
	resources := []*db.Resource{
		{ResourceName: "课件1", FilePath: "/path/1.pdf", FileType: "pdf", FileSize: 1024, UploaderID: 201, CourseID: courseID, Status: "normal", CreatedAt: time.Now()},
		{ResourceName: "课件2", FilePath: "/path/2.docx", FileType: "docx", FileSize: 2048, UploaderID: 202, CourseID: courseID, Status: "normal", CreatedAt: time.Now()},
		{ResourceName: "课件3", FilePath: "/path/3.pdf", FileType: "pdf", FileSize: 3072, UploaderID: 203, CourseID: courseID, Status: "normal", CreatedAt: time.Now()},
	}
	for _, r := range resources {
		if err := db.DB.WithContext(ctx).Table(constants.ResourceTableName).Create(r).Error; err != nil {
			t.Fatalf("插入测试资源失败: %v", err)
		}
	}

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/api/courses/:course_id/resources", GetCourseResourceList)

	w := ut.PerformRequest(router, "GET", "/api/courses/1/resources?type=pdf&page_num=1&page_size=10", nil)
	resp := w.Result()

	if resp.StatusCode() != 200 {
		t.Errorf("期望状态码为 200, 实际为 %d", resp.StatusCode())
	}

	var resourceResp courseModel.GetCourseResourceListResp
	if err := json.Unmarshal(resp.Body(), &resourceResp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if len(resourceResp.Resources) != 2 {
		t.Errorf("期望PDF资源数量为 2, 实际为 %d", len(resourceResp.Resources))
	}
}
