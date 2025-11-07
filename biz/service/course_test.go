package service

import (
	"context"
	"testing"
	"time"

	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/course"
	"LearnShare/pkg/constants"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupCourseTestDB 初始化课程测试数据库
func setupCourseTestDB(t *testing.T) func() {
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

// seedCourse 插入测试课程
func seedCourse(t *testing.T, courseName string, teacherID, majorID int64, grade string) *db.Course {
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

// buildRequestContextWithUID 构建带用户ID的请求上下文
func buildRequestContextWithUID(uid int64) *app.RequestContext {
	ctx := app.NewContext(0)
	ctx.Set(constants.ContextUid, uid)
	return ctx
}

// TestCourseServiceSearch 测试搜索课程
func TestCourseServiceSearch(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	// 插入测试数据
	seedCourse(t, "高等数学", 101, 1, "2024")
	seedCourse(t, "线性代数", 102, 1, "2024")
	seedCourse(t, "概率论", 103, 1, "2023")
	seedCourse(t, "数据结构", 104, 2, "2024")

	svc := NewCourseService(context.Background(), nil)

	// 测试关键词搜索
	keywords := "数学"
	req := &course.SearchReq{
		Keywords: &keywords,
		PageNum:  1,
		PageSize: 10,
	}
	courses, err := svc.Search(req)
	if err != nil {
		t.Fatalf("搜索课程失败: %v", err)
	}
	if len(courses) != 1 {
		t.Errorf("期望搜索到 1 门课程, 实际为 %d", len(courses))
	}
	if courses[0].CourseName != "高等数学" {
		t.Errorf("期望课程名为 高等数学, 实际为 %s", courses[0].CourseName)
	}

	// 测试年级过滤
	grade := "2024"
	req2 := &course.SearchReq{
		Grade:    &grade,
		PageNum:  1,
		PageSize: 10,
	}
	courses2, err := svc.Search(req2)
	if err != nil {
		t.Fatalf("按年级搜索失败: %v", err)
	}
	if len(courses2) != 3 {
		t.Errorf("期望搜索到 3 门课程, 实际为 %d", len(courses2))
	}
}

// TestCourseServiceSearchWithEmptyKeywords 测试空关键词搜索
func TestCourseServiceSearchWithEmptyKeywords(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	seedCourse(t, "计算机网络", 105, 2, "2024")
	seedCourse(t, "操作系统", 106, 2, "2024")

	svc := NewCourseService(context.Background(), nil)
	req := &course.SearchReq{
		PageNum:  1,
		PageSize: 10,
	}
	courses, err := svc.Search(req)
	if err != nil {
		t.Fatalf("搜索失败: %v", err)
	}
	if len(courses) != 2 {
		t.Errorf("期望返回 2 门课程, 实际为 %d", len(courses))
	}
}

// TestCourseServiceSearchPagination 测试分页功能
func TestCourseServiceSearchPagination(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	// 插入多条数据
	for i := 1; i <= 5; i++ {
		seedCourse(t, "课程", int64(100+i), 1, "2024")
	}

	svc := NewCourseService(context.Background(), nil)
	req := &course.SearchReq{
		PageNum:  1,
		PageSize: 2,
	}
	courses, err := svc.Search(req)
	if err != nil {
		t.Fatalf("分页搜索失败: %v", err)
	}
	if len(courses) != 2 {
		t.Errorf("期望返回 2 门课程, 实际为 %d", len(courses))
	}

	// 测试第二页
	req.PageNum = 2
	courses2, err := svc.Search(req)
	if err != nil {
		t.Fatalf("获取第二页失败: %v", err)
	}
	if len(courses2) != 2 {
		t.Errorf("期望第二页返回 2 门课程, 实际为 %d", len(courses2))
	}
}

// TestCourseServiceGetCourseDetail 测试获取课程详情
func TestCourseServiceGetCourseDetail(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	inserted := seedCourse(t, "编译原理", 107, 2, "2024")

	svc := NewCourseService(context.Background(), nil)
	req := &course.GetCourseDetailReq{CourseID: inserted.CourseID}
	courseDetail, err := svc.GetCourseDetail(req)
	if err != nil {
		t.Fatalf("获取课程详情失败: %v", err)
	}
	if courseDetail.CourseName != "编译原理" {
		t.Errorf("期望课程名为 编译原理, 实际为 %s", courseDetail.CourseName)
	}
	if courseDetail.TeacherId != 107 {
		t.Errorf("期望教师ID为 107, 实际为 %d", courseDetail.TeacherId)
	}
}

// TestCourseServiceGetCourseResourceList 测试获取课程资源列表
func TestCourseServiceGetCourseResourceList(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	courseID := int64(501)

	// 插入测试资源
	resources := []*db.Resource{
		{ResourceName: "课件1", FilePath: "/path/1.pdf", FileType: "pdf", FileSize: 1024, UploaderID: 201, CourseID: courseID, Status: "normal", CreatedAt: time.Now()},
		{ResourceName: "课件2", FilePath: "/path/2.docx", FileType: "docx", FileSize: 2048, UploaderID: 202, CourseID: courseID, Status: "normal", CreatedAt: time.Now()},
		{ResourceName: "课件3", FilePath: "/path/3.pdf", FileType: "pdf", FileSize: 3072, UploaderID: 203, CourseID: courseID, Status: "pending_review", CreatedAt: time.Now()},
	}
	for _, r := range resources {
		if err := db.DB.WithContext(ctx).Table(constants.ResourceTableName).Create(r).Error; err != nil {
			t.Fatalf("插入测试资源失败: %v", err)
		}
	}

	svc := NewCourseService(context.Background(), nil)
	req := &course.GetCourseResourceListReq{
		CourseID: courseID,
		PageNum:  1,
		PageSize: 10,
	}
	resourceList, err := svc.GetCourseResourceList(req)
	if err != nil {
		t.Fatalf("获取课程资源列表失败: %v", err)
	}
	if len(resourceList) != 3 {
		t.Errorf("期望资源数量为 3, 实际为 %d", len(resourceList))
	}

	// 测试类型过滤
	fileType := "pdf"
	req2 := &course.GetCourseResourceListReq{
		CourseID: courseID,
		Type:     &fileType,
		PageNum:  1,
		PageSize: 10,
	}
	pdfResources, err := svc.GetCourseResourceList(req2)
	if err != nil {
		t.Fatalf("按类型过滤失败: %v", err)
	}
	if len(pdfResources) != 2 {
		t.Errorf("期望PDF资源数量为 2, 实际为 %d", len(pdfResources))
	}

	// 测试状态过滤
	status := "normal"
	req3 := &course.GetCourseResourceListReq{
		CourseID: courseID,
		Status:   &status,
		PageNum:  1,
		PageSize: 10,
	}
	normalResources, err := svc.GetCourseResourceList(req3)
	if err != nil {
		t.Fatalf("按状态过滤失败: %v", err)
	}
	if len(normalResources) != 2 {
		t.Errorf("期望normal状态资源数量为 2, 实际为 %d", len(normalResources))
	}
}

// TestCourseServiceGetCourseComments 测试获取课程评论列表
func TestCourseServiceGetCourseComments(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	courseID := int64(502)

	// 插入测试评论
	comments := []*db.CourseComment{
		{CourseID: courseID, UserID: 301, Content: "很好的课程", ParentID: 0, IsVisible: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CourseID: courseID, UserID: 302, Content: "老师讲得很清楚", ParentID: 0, IsVisible: true, CreatedAt: time.Now().Add(-time.Hour), UpdatedAt: time.Now()},
		{CourseID: courseID, UserID: 303, Content: "不可见的评论", ParentID: 0, IsVisible: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, c := range comments {
		if err := db.DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(c).Error; err != nil {
			t.Fatalf("插入测试评论失败: %v", err)
		}
	}

	svc := NewCourseService(context.Background(), nil)
	req := &course.GetCourseCommentsReq{
		CourseID: courseID,
		SortBy:   "latest",
		PageNum:  1,
		PageSize: 10,
	}
	commentList, err := svc.GetCourseComments(req)
	if err != nil {
		t.Fatalf("获取评论列表失败: %v", err)
	}
	if len(commentList) != 2 { // 只返回可见评论
		t.Errorf("期望评论数量为 2, 实际为 %d", len(commentList))
	}
}

// TestCourseServiceGetCourseCommentsWithDifferentSorting 测试不同排序方式
func TestCourseServiceGetCourseCommentsWithDifferentSorting(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	courseID := int64(503)

	// 插入测试评论
	comments := []*db.CourseComment{
		{CourseID: courseID, UserID: 304, Content: "最新评论", ParentID: 0, IsVisible: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CourseID: courseID, UserID: 305, Content: "较早评论", ParentID: 0, IsVisible: true, CreatedAt: time.Now().Add(-2 * time.Hour), UpdatedAt: time.Now()},
	}
	for _, c := range comments {
		if err := db.DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(c).Error; err != nil {
			t.Fatalf("插入测试评论失败: %v", err)
		}
	}

	svc := NewCourseService(context.Background(), nil)

	// 测试最新排序
	req1 := &course.GetCourseCommentsReq{
		CourseID: courseID,
		SortBy:   "latest",
		PageNum:  1,
		PageSize: 10,
	}
	latestComments, err := svc.GetCourseComments(req1)
	if err != nil {
		t.Fatalf("获取最新评论失败: %v", err)
	}
	if len(latestComments) > 0 && latestComments[0].Content != "最新评论" {
		t.Errorf("期望第一条评论为 最新评论, 实际为 %s", latestComments[0].Content)
	}

	// 测试最早排序
	req2 := &course.GetCourseCommentsReq{
		CourseID: courseID,
		SortBy:   "oldest",
		PageNum:  1,
		PageSize: 10,
	}
	oldestComments, err := svc.GetCourseComments(req2)
	if err != nil {
		t.Fatalf("获取最早评论失败: %v", err)
	}
	if len(oldestComments) > 0 && oldestComments[0].Content != "较早评论" {
		t.Errorf("期望第一条评论为 较早评论, 实际为 %s", oldestComments[0].Content)
	}
}

// TestCourseServiceSubmitCourseRating 测试提交课程评分
func TestCourseServiceSubmitCourseRating(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	ctx := buildRequestContextWithUID(401)
	svc := NewCourseService(context.Background(), ctx)

	req := &course.SubmitCourseRatingReq{
		CourseID: 601,
		Rating:   5,
	}
	err := svc.SubmitCourseRating(req)
	if err != nil {
		t.Fatalf("提交课程评分失败: %v", err)
	}

	// 验证评分已保存
	var rating db.CourseRating
	if err := db.DB.WithContext(context.Background()).Table(constants.CourseRatingTableName).
		Where("user_id = ? AND course_id = ?", 401, 601).First(&rating).Error; err != nil {
		t.Fatalf("查询评分失败: %v", err)
	}
	if rating.Recommendation != 5 {
		t.Errorf("期望推荐度为 5, 实际为 %d", rating.Recommendation)
	}
}

// TestCourseServiceSubmitCourseComment 测试提交课程评论
func TestCourseServiceSubmitCourseComment(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	ctx := buildRequestContextWithUID(402)
	svc := NewCourseService(context.Background(), ctx)

	req := &course.SubmitCourseCommentReq{
		CourseID:  602,
		Contents:  "这是一个测试评论",
		ParentID:  0,
		IsVisible: true,
	}
	err := svc.SubmitCourseComment(req)
	if err != nil {
		t.Fatalf("提交课程评论失败: %v", err)
	}

	// 验证评论已保存
	var comment db.CourseComment
	if err := db.DB.WithContext(context.Background()).Table(constants.CourseCommentTableName).
		Where("user_id = ? AND course_id = ?", 402, 602).First(&comment).Error; err != nil {
		t.Fatalf("查询评论失败: %v", err)
	}
	if comment.Content != "这是一个测试评论" {
		t.Errorf("期望评论内容为 这是一个测试评论, 实际为 %s", comment.Content)
	}
}

// TestCourseServiceDeleteCourseComment 测试删除课程评论
func TestCourseServiceDeleteCourseComment(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	comment := &db.CourseComment{
		CourseID:  603,
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

	svc := NewCourseService(context.Background(), nil)
	req := &course.DeleteCourseCommentReq{CommentID: comment.CommentID}
	err := svc.DeleteCourseComment(req)
	if err != nil {
		t.Fatalf("删除评论失败: %v", err)
	}

	// 验证评论已删除
	var deleted db.CourseComment
	err = db.DB.WithContext(ctx).Table(constants.CourseCommentTableName).
		Where("comment_id = ?", comment.CommentID).First(&deleted).Error
	if err == nil {
		t.Error("期望查询不到已删除的评论")
	}
}

// TestCourseServiceDeleteCourseRating 测试删除课程评分
func TestCourseServiceDeleteCourseRating(t *testing.T) {
	cleanup := setupCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rating := &db.CourseRating{
		UserID:         404,
		CourseID:       604,
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

	svc := NewCourseService(context.Background(), nil)
	req := &course.DeleteCourseRatingReq{RatingID: rating.RatingID}
	err := svc.DeleteCourseRating(req)
	if err != nil {
		t.Fatalf("删除评分失败: %v", err)
	}

	// 验证评分已删除
	var deleted db.CourseRating
	err = db.DB.WithContext(ctx).Table(constants.CourseRatingTableName).
		Where("rating_id = ?", rating.RatingID).First(&deleted).Error
	if err == nil {
		t.Error("期望查询不到已删除的评分")
	}
}
