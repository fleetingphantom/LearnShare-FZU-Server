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

// CreateCourseAsync 异步创建课程
func CreateCourseAsync(ctx context.Context, courseName string, teacherID, majorID int64, credit float64, grade, description string) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return CreateCourse(ctx, courseName, teacherID, majorID, credit, grade, description)
	})
}

func UpdateCourse(ctx context.Context, courseID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", courseID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程失败: "+err.Error())
	}
	return nil
}

// UpdateCourseAsync 异步更新课程
func UpdateCourseAsync(ctx context.Context, courseID int64, updates map[string]interface{}) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateCourse(ctx, courseID, updates)
	})
}

func DeleteCourse(ctx context.Context, courseID int64) error {
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", courseID).Delete(&Course{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程失败: "+err.Error())
	}
	return nil
}

// DeleteCourseAsync 异步删除课程
func DeleteCourseAsync(ctx context.Context, courseID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteCourse(ctx, courseID)
	})
}

func GetCourseByID(ctx context.Context, courseID int64) (*Course, error) {
	var course Course
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", courseID).First(&course).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程失败: "+err.Error())
	}
	return &course, nil
}

func GetCoursesByTeacherID(ctx context.Context, teacherID int64, pageSize, pageNum int) ([]*Course, error) {
	var courses []*Course
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("teacher_id = ?", teacherID).Limit(pageSize).Offset(pageSize * (pageNum - 1)).Find(&courses).Error
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

func SearchCourses(ctx context.Context, keywords string, grade string, pageNum, pageSize int) ([]*Course, error) {
	var courses []*Course // 声明courses变量

	query := DB.WithContext(ctx).Table(constants.CourseTableName)

	if keywords != "" {
		query = query.Where("course_name LIKE ?", "%"+keywords+"%")
	}
	if grade != "" {
		query = query.Where("grade = ?", grade)
	}

	err := query.Limit(pageSize).Offset(pageSize * (pageNum - 1)).Find(&courses).Error
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

// SubmitCourseRatingAsync 异步提交课程评分
func SubmitCourseRatingAsync(ctx context.Context, rating *CourseRating) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return SubmitCourseRating(ctx, rating)
	})
}

func UpdateCourseRating(ctx context.Context, ratingID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程评分失败: "+err.Error())
	}
	return nil
}

// UpdateCourseRatingAsync 异步更新课程评分
func UpdateCourseRatingAsync(ctx context.Context, ratingID int64, updates map[string]interface{}) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateCourseRating(ctx, ratingID, updates)
	})
}

func DeleteCourseRating(ctx context.Context, ratingID int64) error {
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).Delete(&CourseRating{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评分失败: "+err.Error())
	}
	return nil
}

// DeleteCourseRatingAsync 异步删除课程评分
func DeleteCourseRatingAsync(ctx context.Context, ratingID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteCourseRating(ctx, ratingID)
	})
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

// SubmitCourseCommentAsync 异步提交课程评论
func SubmitCourseCommentAsync(ctx context.Context, comment *CourseComment) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return SubmitCourseComment(ctx, comment)
	})
}

func UpdateCourseComment(ctx context.Context, commentID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程评论失败: "+err.Error())
	}
	return nil
}

// UpdateCourseCommentAsync 异步更新课程评论
func UpdateCourseCommentAsync(ctx context.Context, commentID int64, updates map[string]interface{}) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateCourseComment(ctx, commentID, updates)
	})
}

func DeleteCourseComment(ctx context.Context, commentID int64) error {
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).Delete(&CourseComment{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评论失败: "+err.Error())
	}
	return nil
}

// DeleteCourseCommentAsync 异步删除课程评论
func DeleteCourseCommentAsync(ctx context.Context, commentID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteCourseComment(ctx, commentID)
	})
}

func GetCourseCommentByID(ctx context.Context, commentID int64) (*CourseComment, error) {
	var comment CourseComment
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).First(&comment).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评论失败: "+err.Error())
	}
	return &comment, nil
}

func GetCourseCommentsByCourseID(ctx context.Context, courseID int64, sortBy string, pageNum, pageSize int) ([]*CourseCommentWithuser, error) {
	// 使用联表查询，一次性获取评论与用户信息

	var rows []CommentUserRow

	query := DB.WithContext(ctx).Table(constants.CourseCommentTableName+" AS c").Select(
		"c.comment_id, c.course_id, c.content, c.parent_id, c.is_visible, c.created_at, c.updated_at, "+
			"u.user_id AS u_user_id, u.username AS u_username, u.email AS u_email, u.college_id AS u_college_id, u.major_id AS u_major_id, u.avatar_url AS u_avatar_url, u.reputation_score AS u_reputation_score, u.role_id AS u_role_id, u.status AS u_status, u.created_at AS u_created_at, u.updated_at AS u_updated_at",
	).Joins("LEFT JOIN "+constants.UserTableName+" u ON c.user_id = u.user_id").Where("c.course_id = ? AND c.is_visible = ?", courseID, true)

	// 排序方式
	switch sortBy {
	case "latest":
		query = query.Order("c.created_at DESC")
	case "oldest":
		query = query.Order("c.created_at ASC")
	case "popular":
		query = query.Order("c.created_at DESC") // 简化处理
	default:
		query = query.Order("c.created_at DESC")
	}

	err := query.Limit(pageSize).Offset(pageSize * (pageNum - 1)).Scan(&rows).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评论列表失败: "+err.Error())
	}

	if len(rows) == 0 {
		return []*CourseCommentWithuser{}, nil
	}

	// 将行映射为最终结构
	result := make([]*CourseCommentWithuser, 0, len(rows))
	for _, r := range rows {
		var user User
		if r.UserID == nil {
			user = User{UserID: 0, Username: "未知用户"}
		} else {
			user = User{
				UserID:          *r.UserID,
				Username:        derefString(r.Username),
				Email:           derefString(r.Email),
				CollegeID:       r.CollegeID,
				MajorID:         r.MajorID,
				AvatarURL:       r.AvatarURL,
				ReputationScore: derefInt64(r.ReputationScore),
				RoleID:          derefInt64(r.RoleID),
				Status:          derefString(r.Status),
			}
		}

		cc := &CourseCommentWithuser{
			CommentID: r.CommentID,
			CourseID:  r.CourseID,
			User:      user,
			Content:   r.Content,
			ParentID:  r.ParentID,
			IsVisible: r.IsVisible,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
		result = append(result, cc)
	}

	return result, nil
}

// 辅助函数，用于安全地解引用可空类型
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// GetCourseResources 获取课程资源列表
func GetCourseResources(ctx context.Context, courseID int64, resourceType, status string, pageNum, pageSize int) ([]*Resource, error) {
	var resources []*Resource

	query := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("course_id = ?", courseID)

	if resourceType != "" {
		query = query.Where("type = ?", resourceType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Limit(pageSize).Offset(pageSize * (pageNum - 1)).Find(&resources).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程资源列表失败: "+err.Error())
	}

	return resources, nil
}

// CreateResource 创建资源
func CreateResource(ctx context.Context, resource *Resource) error {
	err := DB.WithContext(ctx).Table(constants.ResourceTableName).Create(resource).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建资源失败: "+err.Error())
	}
	return nil
}

// CreateResourceAsync 异步创建资源
func CreateResourceAsync(ctx context.Context, resource *Resource) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return CreateResource(ctx, resource)
	})
}

// UpdateResource 更新资源
func UpdateResource(ctx context.Context, resourceID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("resource_id = ?", resourceID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源失败: "+err.Error())
	}
	return nil
}

// UpdateResourceAsync 异步更新资源
func UpdateResourceAsync(ctx context.Context, resourceID int64, updates map[string]interface{}) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateResource(ctx, resourceID, updates)
	})
}

// 删除资源
func DeleteResource(ctx context.Context, resourceID int64) error {
	err := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("resource_id = ?", resourceID).Delete(&Resource{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除资源失败: "+err.Error())
	}
	return nil
}

// DeleteResourceAsync 异步删除资源
func DeleteResourceAsync(ctx context.Context, resourceID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteResource(ctx, resourceID)
	})
}
