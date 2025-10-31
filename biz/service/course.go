package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/course"
	"LearnShare/biz/model/module"
	"LearnShare/pkg/errno"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type CourseService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewCourseService(ctx context.Context, c *app.RequestContext) *CourseService {
	return &CourseService{ctx: ctx, c: c}
}

// 搜索课程（对应 Thrift Search 接口）
func (s *CourseService) Search(req *course.SearchReq) ([]*module.Course, error) {
	courses, err := db.SearchCourses(s.ctx, req)
	if err != nil {
		return nil, err
	}

	var moduleCourses []*module.Course
	for _, dbCourse := range courses {
		moduleCourses = append(moduleCourses, dbCourse.ToCourseModule())
	}
	return moduleCourses, nil
}

// 获取课程详情（对应 Thrift GetCourseDetail 接口）
func (s *CourseService) GetCourseDetail(req *course.GetCourseDetailReq) (*module.Course, error) {
	courseInfo, err := db.GetCourseByID(s.ctx, req.CourseId)
	if err != nil {
		return nil, err
	}
	return courseInfo.ToCourseModule(), nil
}

// 获取课程资源列表（对应 Thrift GetCourseResourceList 接口）
func (s *CourseService) GetCourseResourceList(req *course.GetCourseResourceListReq) ([]*module.Resource, error) {
	resources, err := db.GetCourseResources(s.ctx, req)
	if err != nil {
		return nil, err
	}

	var moduleResources []*module.Resource
	for _, dbRes := range resources {
		moduleResources = append(moduleResources, dbRes.ToResourceModule())
	}
	return moduleResources, nil
}

// 获取课程评论列表（对应 Thrift GetCourseComments 接口）
func (s *CourseService) GetCourseComments(req *course.GetCourseCommentsReq) ([]*module.CourseComment, error) {
	comments, err := db.GetCourseComments(s.ctx, req)
	if err != nil {
		return nil, err
	}

	var moduleComments []*module.CourseComment
	for _, dbComment := range comments {
		moduleComments = append(moduleComments, dbComment.ToCourseCommentModule())
	}
	return moduleComments, nil
}

// 提交课程评分（对应 Thrift SubmitCourseRating 接口）
func (s *CourseService) SubmitCourseRating(req *course.SubmitCourseRatingReq) (*module.CourseRating, error) {
	userId := s.getUidFromContext()
	if userId == 0 {
		return nil, errno.NewErrNo(errno.AuthErrCode, "未登录")
	}

	if req.UserId != userId {
		return nil, errno.NewErrNo(errno.PermissionErrCode, "无提交权限")
	}

	err := db.CreateCourseRating(s.ctx, req)
	if err != nil {
		return nil, err
	}

	rating, err := db.GetCourseRatingByID(s.ctx, req.RatingId)
	if err != nil {
		return nil, err
	}
	return rating.ToCourseRatingModule(), nil
}

// 提交课程评论（对应 Thrift SubmitCourseComment 接口）
func (s *CourseService) SubmitCourseComment(req *course.SubmitCourseCommentReq) (*module.CourseComment, error) {
	userId := s.getUidFromContext()
	if userId == 0 {
		return nil, errno.NewErrNo(errno.AuthErrCode, "未登录")
	}

	comment, err := db.CreateCourseComment(s.ctx, userId, req)
	if err != nil {
		return nil, err
	}
	return comment.ToCourseCommentModule(), nil
}

// 删除课程评论（对应 Thrift DeleteCourseComment 接口）
func (s *CourseService) DeleteCourseComment(req *course.DeleteCourseCommentReq) error {
	userId := s.getUidFromContext()
	if userId == 0 {
		return errno.NewErrNo(errno.AuthErrCode, "未登录")
	}

	hasPermission, err := db.CheckCommentOwner(s.ctx, req.CommentId, userId)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errno.NewErrNo(errno.PermissionErrCode, "无删除权限")
	}

	return db.DeleteCourseComment(s.ctx, req.CommentId)
}

// 删除课程评分（对应 Thrift DeleteCourseRating 接口）
func (s *CourseService) DeleteCourseRating(req *course.DeleteCourseRatingReq) error {
	userId := s.getUidFromContext()
	if userId == 0 {
		return errno.NewErrNo(errno.AuthErrCode, "未登录")
	}

	hasPermission, err := db.CheckRatingOwner(s.ctx, req.RatingId, userId)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errno.NewErrNo(errno.PermissionErrCode, "无删除权限")
	}

	return db.DeleteCourseRating(s.ctx, req.RatingId)
}

// 内部方法：获取当前用户ID（避免与user.go重复声明）
func (s *CourseService) getUidFromContext() int64 {
	uid, ok := s.c.Get("user_id")
	if !ok {
		return 0
	}
	uidInt64, ok := uid.(int64)
	if !ok {
		return 0
	}
	return uidInt64
}
