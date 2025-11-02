package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
)

func CreateCourse(ctx context.Context, courseName string, teacherID, majorID int64, credit float64, grade, description string) error {
	course := &Course{
		CourseName:  courseName,
		TeacherID:   teacherID,
		Credit:      credit,
		MajorID:     majorID,
		Grade:       grade,
		Description: &description,
	}

	err := DB.WithContext(ctx).Table(constants.CourseTableName).Create(course).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建课程失败: "+err.Error())
	}
	return nil
}

func UpdateCourse(ctx context.Context, courseID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", courseID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程失败: "+err.Error())
	}
	return nil
}

func DeleteCourse(ctx context.Context, courseID int64) error {
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", courseID).Delete(&Course{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程失败: "+err.Error())
	}
	return nil
}

func GetCourseByID(ctx context.Context, courseID int64) (*Course, error) {
	var course Course
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", courseID).First(&course).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程失败: "+err.Error())
	}
	return &course, nil
}

func GetCoursesByTeacherID(ctx context.Context, teacherID int64) ([]*Course, error) {
	var courses []*Course
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("teacher_id = ?", teacherID).Find(&courses).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询教师课程失败: "+err.Error())
	}
	return courses, nil
}

func GetCoursesByMajorID(ctx context.Context, majorID int64) ([]*Course, error) {
	var courses []*Course
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("major_id = ?", majorID).Find(&courses).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询专业课程失败: "+err.Error())
	}
	return courses, nil
}

func SearchCourses(ctx context.Context, keywords string, grade string) ([]*Course, error) {
	var courses []*Course // 声明courses变量

	query := DB.WithContext(ctx).Table(constants.CourseTableName)

	if keywords != "" {
		query = query.Where("course_name LIKE ?", "%"+keywords+"%")
	}
	if grade != "" {
		query = query.Where("grade = ?", grade)
	}

	err := query.Find(&courses).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "搜索课程失败: "+err.Error())
	}
	return courses, nil
}

func SubmitCourseRating(ctx context.Context, rating *CourseRating) error {
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Create(rating).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交课程评分失败: "+err.Error())
	}
	return nil
}

func UpdateCourseRating(ctx context.Context, ratingID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程评分失败: "+err.Error())
	}
	return nil
}

func DeleteCourseRating(ctx context.Context, ratingID int64) error {
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).Delete(&CourseRating{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评分失败: "+err.Error())
	}
	return nil
}

func GetCourseRatingByID(ctx context.Context, ratingID int64) (*CourseRating, error) {
	var rating CourseRating
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).First(&rating).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评分失败: "+err.Error())
	}
	return &rating, nil
}

func GetCourseRatingsByCourseID(ctx context.Context, courseID int64) ([]*CourseRating, error) {
	var ratings []*CourseRating
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("course_id = ? AND is_visible = ?", courseID, true).Order("created_at DESC").Find(&ratings).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评分列表失败: "+err.Error())
	}
	return ratings, nil
}

func SubmitCourseComment(ctx context.Context, comment *CourseComment) error {
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(comment).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交课程评论失败: "+err.Error())
	}
	return nil
}

func UpdateCourseComment(ctx context.Context, commentID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程评论失败: "+err.Error())
	}
	return nil
}

func DeleteCourseComment(ctx context.Context, commentID int64) error {
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).Delete(&CourseComment{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评论失败: "+err.Error())
	}
	return nil
}

func GetCourseCommentByID(ctx context.Context, commentID int64) (*CourseComment, error) {
	var comment CourseComment
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).First(&comment).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评论失败: "+err.Error())
	}
	return &comment, nil
}

func GetCourseCommentsByCourseID(ctx context.Context, courseID int64, sortBy string) ([]*CourseComment, error) {
	var comments []*CourseComment

	query := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("course_id = ? AND is_visible = ?", courseID, true)

	// 排序方式
	switch sortBy {
	case "latest":
		query = query.Order("created_at DESC")
	case "oldest":
		query = query.Order("created_at ASC")
	case "popular":
		query = query.Order("created_at DESC") // 简化处理
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Find(&comments).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评论列表失败: "+err.Error())
	}

	return comments, nil
}

// 获取课程资源列表
func GetCourseResources(ctx context.Context, courseID int64, resourceType, status string) ([]*Resource, error) {
	var resources []*Resource

	query := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("course_id = ?", courseID)

	if resourceType != "" {
		query = query.Where("type = ?", resourceType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Find(&resources).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程资源列表失败: "+err.Error())
	}

	return resources, nil
}

// 创建资源
func CreateResource(ctx context.Context, resource *Resource) error {
	err := DB.WithContext(ctx).Table(constants.ResourceTableName).Create(resource).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建资源失败: "+err.Error())
	}
	return nil
}

// 更新资源
func UpdateResource(ctx context.Context, resourceID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("resource_id = ?", resourceID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源失败: "+err.Error())
	}
	return nil
}

// 删除资源
func DeleteResource(ctx context.Context, resourceID int64) error {
	err := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("resource_id = ?", resourceID).Delete(&Resource{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除资源失败: "+err.Error())
	}
	return nil
}
