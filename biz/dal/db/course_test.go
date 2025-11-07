package db

import (
	"context"
	"testing"
	"time"

	"LearnShare/pkg/constants"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// initCourseTestDB 初始化课程测试数据库
func initCourseTestDB(t *testing.T) func() {
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

	DB = sqliteDB

	return func() {
		sqlDB, err := DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

// 辅助函数：插入测试课程
func insertTestCourse(t *testing.T, courseName string, teacherID, majorID int64) *Course {
	t.Helper()
	now := time.Now()
	description := "测试课程描述"
	course := &Course{
		CourseName:  courseName,
		TeacherID:   teacherID,
		Credit:      3.0,
		MajorID:     majorID,
		Grade:       "2024",
		Description: &description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := DB.WithContext(context.Background()).Table(constants.CourseTableName).Create(course).Error; err != nil {
		t.Fatalf("插入测试课程失败: %v", err)
	}
	return course
}

// TestCreateCourse 测试创建课程
func TestCreateCourse(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	err := CreateCourse(ctx, "高等数学", 101, 1, 4.0, "2024", "数学基础课程")
	if err != nil {
		t.Fatalf("创建课程失败: %v", err)
	}

	var course Course
	if err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_name = ?", "高等数学").First(&course).Error; err != nil {
		t.Fatalf("查询课程失败: %v", err)
	}

	if course.CourseName != "高等数学" {
		t.Errorf("期望课程名为 高等数学, 实际为 %s", course.CourseName)
	}
	if course.TeacherID != 101 {
		t.Errorf("期望教师ID为 101, 实际为 %d", course.TeacherID)
	}
	if course.Credit != 4.0 {
		t.Errorf("期望学分为 4.0, 实际为 %f", course.Credit)
	}
}

// TestUpdateCourse 测试更新课程
func TestUpdateCourse(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	course := insertTestCourse(t, "线性代数", 102, 1)
	ctx := context.Background()

	updates := map[string]interface{}{
		"course_name": "线性代数A",
		"credit":      3.5,
	}

	err := UpdateCourse(ctx, course.CourseID, updates)
	if err != nil {
		t.Fatalf("更新课程失败: %v", err)
	}

	var updated Course
	if err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", course.CourseID).First(&updated).Error; err != nil {
		t.Fatalf("查询更新后的课程失败: %v", err)
	}

	if updated.CourseName != "线性代数A" {
		t.Errorf("期望课程名为 线性代数A, 实际为 %s", updated.CourseName)
	}
	if updated.Credit != 3.5 {
		t.Errorf("期望学分为 3.5, 实际为 %f", updated.Credit)
	}
}

// TestDeleteCourse 测试删除课程
func TestDeleteCourse(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	course := insertTestCourse(t, "概率论", 103, 1)
	ctx := context.Background()

	err := DeleteCourse(ctx, course.CourseID)
	if err != nil {
		t.Fatalf("删除课程失败: %v", err)
	}

	var deleted Course
	err = DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", course.CourseID).First(&deleted).Error
	if err == nil {
		t.Error("期望查询不到已删除的课程")
	}
}

// TestGetCourseByID 测试根据ID获取课程
func TestGetCourseByID(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	inserted := insertTestCourse(t, "数据结构", 104, 2)
	ctx := context.Background()

	course, err := GetCourseByID(ctx, inserted.CourseID)
	if err != nil {
		t.Fatalf("获取课程失败: %v", err)
	}

	if course.CourseName != "数据结构" {
		t.Errorf("期望课程名为 数据结构, 实际为 %s", course.CourseName)
	}
	if course.TeacherID != 104 {
		t.Errorf("期望教师ID为 104, 实际为 %d", course.TeacherID)
	}
}

// TestGetCoursesByTeacherID 测试获取教师课程列表
func TestGetCoursesByTeacherID(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	teacherID := int64(105)

	// 插入多个课程
	insertTestCourse(t, "计算机网络", teacherID, 2)
	insertTestCourse(t, "操作系统", teacherID, 2)
	insertTestCourse(t, "编译原理", 999, 2) // 不同教师

	courses, err := GetCoursesByTeacherID(ctx, teacherID, 10, 1)
	if err != nil {
		t.Fatalf("获取教师课程列表失败: %v", err)
	}

	if len(courses) != 2 {
		t.Errorf("期望课程数量为 2, 实际为 %d", len(courses))
	}
}

// TestGetCoursesByMajorID 测试获取专业课程列表
func TestGetCoursesByMajorID(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	majorID := int64(3)

	insertTestCourse(t, "微观经济学", 106, majorID)
	insertTestCourse(t, "宏观经济学", 107, majorID)
	insertTestCourse(t, "计量经济学", 108, 999) // 不同专业

	courses, err := GetCoursesByMajorID(ctx, majorID)
	if err != nil {
		t.Fatalf("获取专业课程列表失败: %v", err)
	}

	if len(courses) != 2 {
		t.Errorf("期望课程数量为 2, 实际为 %d", len(courses))
	}
}

// TestSearchCourses 测试搜索课程
func TestSearchCourses(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()

	insertTestCourse(t, "大学物理", 109, 1)
	insertTestCourse(t, "大学英语", 110, 1)
	insertTestCourse(t, "高等物理", 111, 1)

	// 测试关键词搜索
	courses, err := SearchCourses(ctx, "物理", "", 1, 10)
	if err != nil {
		t.Fatalf("搜索课程失败: %v", err)
	}

	if len(courses) != 2 {
		t.Errorf("期望搜索到 2 门课程, 实际为 %d", len(courses))
	}
}

// TestSubmitCourseRating 测试提交课程评分
func TestSubmitCourseRating(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rating := &CourseRating{
		UserID:         201,
		CourseID:       301,
		Recommendation: 5,
		Difficulty:     "medium",
		Workload:       3,
		Usefulness:     4,
		IsVisible:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := SubmitCourseRating(ctx, rating)
	if err != nil {
		t.Fatalf("提交课程评分失败: %v", err)
	}

	var saved CourseRating
	if err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("user_id = ? AND course_id = ?", 201, 301).First(&saved).Error; err != nil {
		t.Fatalf("查询评分失败: %v", err)
	}

	if saved.Recommendation != 5 {
		t.Errorf("期望推荐度为 5, 实际为 %d", saved.Recommendation)
	}
}

// TestUpdateCourseRating 测试更新课程评分
func TestUpdateCourseRating(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rating := &CourseRating{
		UserID:         202,
		CourseID:       302,
		Recommendation: 4,
		Difficulty:     "easy",
		Workload:       2,
		Usefulness:     3,
		IsVisible:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Create(rating).Error; err != nil {
		t.Fatalf("插入测试评分失败: %v", err)
	}

	updates := map[string]interface{}{
		"recommendation": 5,
		"difficulty":     "hard",
	}

	err := UpdateCourseRating(ctx, rating.RatingID, updates)
	if err != nil {
		t.Fatalf("更新评分失败: %v", err)
	}

	var updated CourseRating
	if err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", rating.RatingID).First(&updated).Error; err != nil {
		t.Fatalf("查询更新后的评分失败: %v", err)
	}

	if updated.Recommendation != 5 {
		t.Errorf("期望推荐度为 5, 实际为 %d", updated.Recommendation)
	}
	if updated.Difficulty != "hard" {
		t.Errorf("期望难度为 hard, 实际为 %s", updated.Difficulty)
	}
}

// TestDeleteCourseRating 测试删除课程评分
func TestDeleteCourseRating(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rating := &CourseRating{
		UserID:         203,
		CourseID:       303,
		Recommendation: 3,
		Difficulty:     "medium",
		Workload:       3,
		Usefulness:     3,
		IsVisible:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Create(rating).Error; err != nil {
		t.Fatalf("插入测试评分失败: %v", err)
	}

	err := DeleteCourseRating(ctx, rating.RatingID)
	if err != nil {
		t.Fatalf("删除评分失败: %v", err)
	}

	var deleted CourseRating
	err = DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", rating.RatingID).First(&deleted).Error
	if err == nil {
		t.Error("期望查询不到已删除的评分")
	}
}

// TestGetCourseRatingByID 测试根据ID获取评分
func TestGetCourseRatingByID(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rating := &CourseRating{
		UserID:         204,
		CourseID:       304,
		Recommendation: 4,
		Difficulty:     "medium",
		Workload:       3,
		Usefulness:     4,
		IsVisible:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Create(rating).Error; err != nil {
		t.Fatalf("插入测试评分失败: %v", err)
	}

	fetched, err := GetCourseRatingByID(ctx, rating.RatingID)
	if err != nil {
		t.Fatalf("获取评分失败: %v", err)
	}

	if fetched.Recommendation != 4 {
		t.Errorf("期望推荐度为 4, 实际为 %d", fetched.Recommendation)
	}
}

// TestGetCourseRatingsByCourseID 测试获取课程评分列表
func TestGetCourseRatingsByCourseID(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	courseID := int64(305)

	// 插入多个评分
	ratings := []*CourseRating{
		{UserID: 205, CourseID: courseID, Recommendation: 5, Difficulty: "easy", Workload: 2, Usefulness: 5, IsVisible: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{UserID: 206, CourseID: courseID, Recommendation: 4, Difficulty: "medium", Workload: 3, Usefulness: 4, IsVisible: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{UserID: 207, CourseID: courseID, Recommendation: 3, Difficulty: "hard", Workload: 4, Usefulness: 3, IsVisible: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}, // 不可见
	}

	for _, r := range ratings {
		if err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Create(r).Error; err != nil {
			t.Fatalf("插入测试评分失败: %v", err)
		}
	}

	fetchedRatings, err := GetCourseRatingsByCourseID(ctx, courseID)
	if err != nil {
		t.Fatalf("获取课程评分列表失败: %v", err)
	}

	if len(fetchedRatings) != 2 { // 只返回可见的评分
		t.Errorf("期望评分数量为 2, 实际为 %d", len(fetchedRatings))
	}
}

// TestSubmitCourseComment 测试提交课程评论
func TestSubmitCourseComment(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	comment := &CourseComment{
		CourseID:  401,
		UserID:    501,
		Content:   "这门课很棒！",
		ParentID:  0,
		IsVisible: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := SubmitCourseComment(ctx, comment)
	if err != nil {
		t.Fatalf("提交评论失败: %v", err)
	}

	var saved CourseComment
	if err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("course_id = ? AND user_id = ?", 401, 501).First(&saved).Error; err != nil {
		t.Fatalf("查询评论失败: %v", err)
	}

	if saved.Content != "这门课很棒！" {
		t.Errorf("期望评论内容为 这门课很棒！, 实际为 %s", saved.Content)
	}
}

// TestUpdateCourseComment 测试更新评论
func TestUpdateCourseComment(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	comment := &CourseComment{
		CourseID:  402,
		UserID:    502,
		Content:   "原始评论",
		ParentID:  0,
		IsVisible: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(comment).Error; err != nil {
		t.Fatalf("插入测试评论失败: %v", err)
	}

	updates := map[string]interface{}{
		"content": "修改后的评论",
	}

	err := UpdateCourseComment(ctx, comment.CommentID, updates)
	if err != nil {
		t.Fatalf("更新评论失败: %v", err)
	}

	var updated CourseComment
	if err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", comment.CommentID).First(&updated).Error; err != nil {
		t.Fatalf("查询更新后的评论失败: %v", err)
	}

	if updated.Content != "修改后的评论" {
		t.Errorf("期望评论内容为 修改后的评论, 实际为 %s", updated.Content)
	}
}

// TestDeleteCourseComment 测试删除评论
func TestDeleteCourseComment(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	comment := &CourseComment{
		CourseID:  403,
		UserID:    503,
		Content:   "待删除的评论",
		ParentID:  0,
		IsVisible: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(comment).Error; err != nil {
		t.Fatalf("插入测试评论失败: %v", err)
	}

	err := DeleteCourseComment(ctx, comment.CommentID)
	if err != nil {
		t.Fatalf("删除评论失败: %v", err)
	}

	var deleted CourseComment
	err = DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", comment.CommentID).First(&deleted).Error
	if err == nil {
		t.Error("期望查询不到已删除的评论")
	}
}

// TestGetCourseCommentByID 测试根据ID获取评论
func TestGetCourseCommentByID(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	comment := &CourseComment{
		CourseID:  404,
		UserID:    504,
		Content:   "测试评论",
		ParentID:  0,
		IsVisible: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(comment).Error; err != nil {
		t.Fatalf("插入测试评论失败: %v", err)
	}

	fetched, err := GetCourseCommentByID(ctx, comment.CommentID)
	if err != nil {
		t.Fatalf("获取评论失败: %v", err)
	}

	if fetched.Content != "测试评论" {
		t.Errorf("期望评论内容为 测试评论, 实际为 %s", fetched.Content)
	}
}

// TestGetCourseCommentsByCourseID 测试获取课程评论列表
func TestGetCourseCommentsByCourseID(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	courseID := int64(405)

	// 插入多个评论
	comments := []*CourseComment{
		{CourseID: courseID, UserID: 505, Content: "评论1", ParentID: 0, IsVisible: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CourseID: courseID, UserID: 506, Content: "评论2", ParentID: 0, IsVisible: true, CreatedAt: time.Now().Add(-time.Hour), UpdatedAt: time.Now()},
		{CourseID: courseID, UserID: 507, Content: "评论3", ParentID: 0, IsVisible: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}, // 不可见
	}

	for _, c := range comments {
		if err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(c).Error; err != nil {
			t.Fatalf("插入测试评论失败: %v", err)
		}
	}

	fetchedComments, err := GetCourseCommentsByCourseID(ctx, courseID, "latest", 1, 10)
	if err != nil {
		t.Fatalf("获取课程评论列表失败: %v", err)
	}

	if len(fetchedComments) != 2 { // 只返回可见的评论
		t.Errorf("期望评论数量为 2, 实际为 %d", len(fetchedComments))
	}
}

// TestGetCourseResources 测试获取课程资源列表
func TestGetCourseResources(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	courseID := int64(406)

	// 插入多个资源
	resources := []*Resource{
		{ResourceName: "课件1", FilePath: "/path/1.pdf", FileType: "pdf", FileSize: 1024, UploaderID: 601, CourseID: courseID, Status: "normal", CreatedAt: time.Now()},
		{ResourceName: "课件2", FilePath: "/path/2.docx", FileType: "docx", FileSize: 2048, UploaderID: 602, CourseID: courseID, Status: "normal", CreatedAt: time.Now()},
		{ResourceName: "课件3", FilePath: "/path/3.pdf", FileType: "pdf", FileSize: 3072, UploaderID: 603, CourseID: 999, Status: "normal", CreatedAt: time.Now()}, // 不同课程
	}

	for _, r := range resources {
		if err := DB.WithContext(ctx).Table(constants.ResourceTableName).Create(r).Error; err != nil {
			t.Fatalf("插入测试资源失败: %v", err)
		}
	}

	fetchedResources, err := GetCourseResources(ctx, courseID, "", "", 1, 10)
	if err != nil {
		t.Fatalf("获取课程资源列表失败: %v", err)
	}

	if len(fetchedResources) != 2 {
		t.Errorf("期望资源数量为 2, 实际为 %d", len(fetchedResources))
	}

	// 测试类型过滤
	pdfResources, err := GetCourseResources(ctx, courseID, "pdf", "", 1, 10)
	if err != nil {
		t.Fatalf("获取PDF资源失败: %v", err)
	}

	if len(pdfResources) != 1 {
		t.Errorf("期望PDF资源数量为 1, 实际为 %d", len(pdfResources))
	}
}

// TestCreateResource 测试创建资源
func TestCreateResource(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	resource := &Resource{
		ResourceName: "新资源",
		FilePath:     "/path/new.pdf",
		FileType:     "pdf",
		FileSize:     4096,
		UploaderID:   604,
		CourseID:     407,
		Status:       "pending_review",
		CreatedAt:    time.Now(),
	}

	err := CreateResource(ctx, resource)
	if err != nil {
		t.Fatalf("创建资源失败: %v", err)
	}

	var saved Resource
	if err := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("resource_name = ?", "新资源").First(&saved).Error; err != nil {
		t.Fatalf("查询资源失败: %v", err)
	}

	if saved.ResourceName != "新资源" {
		t.Errorf("期望资源名为 新资源, 实际为 %s", saved.ResourceName)
	}
}

// TestUpdateResource 测试更新资源
func TestUpdateResource(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	resource := &Resource{
		ResourceName: "旧资源",
		FilePath:     "/path/old.pdf",
		FileType:     "pdf",
		FileSize:     5120,
		UploaderID:   605,
		CourseID:     408,
		Status:       "pending_review",
		CreatedAt:    time.Now(),
	}

	if err := DB.WithContext(ctx).Table(constants.ResourceTableName).Create(resource).Error; err != nil {
		t.Fatalf("插入测试资源失败: %v", err)
	}

	updates := map[string]interface{}{
		"resource_name": "新资源名",
		"status":        "normal",
	}

	err := UpdateResource(ctx, resource.ResourceID, updates)
	if err != nil {
		t.Fatalf("更新资源失败: %v", err)
	}

	var updated Resource
	if err := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("resource_id = ?", resource.ResourceID).First(&updated).Error; err != nil {
		t.Fatalf("查询更新后的资源失败: %v", err)
	}

	if updated.ResourceName != "新资源名" {
		t.Errorf("期望资源名为 新资源名, 实际为 %s", updated.ResourceName)
	}
	if updated.Status != "normal" {
		t.Errorf("期望状态为 normal, 实际为 %s", updated.Status)
	}
}

// TestDeleteResource 测试删除资源
func TestDeleteResource(t *testing.T) {
	cleanup := initCourseTestDB(t)
	defer cleanup()

	ctx := context.Background()
	resource := &Resource{
		ResourceName: "待删除资源",
		FilePath:     "/path/delete.pdf",
		FileType:     "pdf",
		FileSize:     6144,
		UploaderID:   606,
		CourseID:     409,
		Status:       "normal",
		CreatedAt:    time.Now(),
	}

	if err := DB.WithContext(ctx).Table(constants.ResourceTableName).Create(resource).Error; err != nil {
		t.Fatalf("插入测试资源失败: %v", err)
	}

	err := DeleteResource(ctx, resource.ResourceID)
	if err != nil {
		t.Fatalf("删除资源失败: %v", err)
	}

	var deleted Resource
	err = DB.WithContext(ctx).Table(constants.ResourceTableName).Where("resource_id = ?", resource.ResourceID).First(&deleted).Error
	if err == nil {
		t.Error("期望查询不到已删除的资源")
	}
}
