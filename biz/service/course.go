package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/course"
	"LearnShare/biz/model/module"
	"LearnShare/pkg/errno"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

type CourseService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewCourseService(ctx context.Context, c *app.RequestContext) *CourseService {
	return &CourseService{ctx: ctx, c: c}
}

func (s *CourseService) Search(req *course.SearchReq) (*course.SearchResp, error) {

	// 使用strings包处理指针类型的参数
	keywords := ""
	if req.Keywords != nil {
		keywords = strings.TrimSpace(*req.Keywords)
	}

	grade := ""
	if req.Grade != nil {
		grade = strings.TrimSpace(*req.Grade)
	}

	// 调用数据库查询课程
	courses, err := db.SearchCourses(s.ctx, keywords, grade)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "搜索课程失败: "+err.Error())
	}

	// 转换为module.Course列表
	var courseModules []*module.Course
	for _, c := range courses {
		courseModules = append(courseModules, c.ToCourseModule())
	}

	// 构建响应
	resp := &course.SearchResp{
		BaseResponse: &module.BaseResp{
			Code:    errno.SuccessCode,
			Message: errno.SuccessMsg,
		},
		Courses: courseModules,
	}

	return resp, nil
}

func (s *CourseService) GetCourseDetail(req *course.GetCourseDetailReq) (*course.GetCourseDetailResp, error) {
	// 获取课程详情
	courseDetail, err := db.GetCourseByID(s.ctx, req.CourseID)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取课程详情失败: "+err.Error())
	}

	resp := &course.GetCourseDetailResp{
		BaseResponse: &module.BaseResp{
			Code:    errno.SuccessCode,
			Message: errno.SuccessMsg,
		},
		Course: courseDetail.ToCourseModule(),
	}
	return resp, nil
}

func (s *CourseService) GetCourseResourceList(req *course.GetCourseResourceListReq) (*course.GetCourseResourceListResp, error) {
	// 处理指针类型的参数
	var resourceType string
	if req.Type != nil {
		resourceType = *req.Type
	}

	var status string
	if req.Status != nil {
		status = *req.Status
	}

	// 获取课程资源列表 - 使用正确的字段名
	resources, err := db.GetCourseResources(s.ctx, req.CourseID, resourceType, status) // 改为 CourseID
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取课程资源失败: "+err.Error())
	}

	// 转换为module.Resource列表
	var resourceModules []*module.Resource
	for _, r := range resources {
		resourceModules = append(resourceModules, r.ToResourceModule())
	}

	resp := &course.GetCourseResourceListResp{
		BaseResponse: &module.BaseResp{
			Code:    errno.SuccessCode,
			Message: errno.SuccessMsg,
		},
		Resources: resourceModules,
	}
	return resp, nil
}

func (s *CourseService) GetCourseComments(req *course.GetCourseCommentsReq) (*course.GetCourseCommentsResp, error) {
	// SortBy 是普通 string 类型，不是指针
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "latest" // 使用默认值
	}

	// 获取课程评论列表
	comments, err := db.GetCourseCommentsByCourseID(s.ctx, req.CourseID, sortBy)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取课程评论失败: "+err.Error())
	}

	// 转换为module.CourseComment列表
	var commentModules []*module.CourseComment
	for _, c := range comments {
		commentModules = append(commentModules, c.ToCourseCommentModule())
	}

	resp := &course.GetCourseCommentsResp{
		BaseResponse: &module.BaseResp{
			Code:    errno.SuccessCode,
			Message: errno.SuccessMsg,
		},
		Comments: commentModules,
	}
	return resp, nil
}

func (s *CourseService) SubmitCourseRating(req *course.SubmitCourseRatingReq) (*course.SubmitCourseRatingResp, error) {
	// 获取用户ID
	userID := GetUidFormContext(s.c)

	// 创建评分对象
	rating := &db.CourseRating{
		UserID:         userID,
		CourseID:       req.CourseID,
		Recommendation: int64(req.Rating), // 直接使用传入的评分
		Difficulty:     "medium",          // 默认值
		Workload:       3,                 // 默认值
		Usefulness:     4,                 // 默认值
		IsVisible:      true,
	}

	// 提交评分
	err := db.SubmitCourseRating(s.ctx, rating)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评分失败: "+err.Error())
	}

	resp := &course.SubmitCourseRatingResp{
		BaseResponse: &module.BaseResp{
			Code:    errno.SuccessCode,
			Message: errno.SuccessMsg,
		},
		Rating: rating.ToCourseRatingModule(),
	}
	return resp, nil
}

func (s *CourseService) SubmitCourseComment(req *course.SubmitCourseCommentReq) (*course.SubmitCourseCommentResp, error) {
	// 获取用户ID
	userID := GetUidFormContext(s.c)

	// 处理 ParentID 默认值
	parentID := req.ParentID


	// 处理 IsVisible 默认值
	isVisible := req.IsVisible
	if !isVisible {
		isVisible = true // 使用默认值
	}

	// 创建评论对象
	comment := &db.CourseComment{
		CourseID:  req.CourseID,
		UserID:    userID,
		Content:   req.Contents,
		ParentID:  parentID,
		IsVisible: isVisible,
	}

	// 提交评论
	err := db.SubmitCourseComment(s.ctx, comment)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评论失败: "+err.Error())
	}

	resp := &course.SubmitCourseCommentResp{
		BaseResponse: &module.BaseResp{
			Code:    errno.SuccessCode,
			Message: errno.SuccessMsg,
		},
		Comment: comment.ToCourseCommentModule(),
	}
	return resp, nil
}

func (s *CourseService) DeleteCourseComment(req *course.DeleteCourseCommentReq) (*course.DeleteCourseCommentResp, error) {
	// 删除评论
	err := db.DeleteCourseComment(s.ctx, req.CommentID)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评论失败: "+err.Error())
	}

	resp := &course.DeleteCourseCommentResp{
		BaseResponse: &module.BaseResp{
			Code:    errno.SuccessCode,
			Message: errno.SuccessMsg,
		},
	}
	return resp, nil
}

func (s *CourseService) DeleteCourseRating(req *course.DeleteCourseRatingReq) (*course.DeleteCourseRatingResp, error) {
	// 删除评分
	err := db.DeleteCourseRating(s.ctx, req.RatingID)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评分失败: "+err.Error())
	}

	resp := &course.DeleteCourseRatingResp{
		BaseResponse: &module.BaseResp{
			Code:    errno.SuccessCode,
			Message: errno.SuccessMsg,
		},
	}
	return resp, nil
}
