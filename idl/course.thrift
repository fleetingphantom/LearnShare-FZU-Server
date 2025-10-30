namespace go course
include "model.thrift"


// 搜索课程
struct SearchReq{
  required i32 page_size
  required i32 page_num
}

struct SearchResp {
  required model.BaseResp baseResponse;
  optional list<model.Course> courses; 
}

// 获取课程详情
struct GetCourseDetailReq {
  required i64 course_id  
}

struct GetCourseDetailResp {
  required model.BaseResp baseResponse;
  optional model.Course course;  
}

// 获取课程资源列表
struct GetCourseResourceListReq {
  required i64 course_id  
  required i32 page_num   
  required i32 page_size  
}

struct GetCourseResourceListResp {
  required model.BaseResp baseResponse;
  optional list<model.Resource> resources; 
}


service CourseService {
  SearchResp search(1: SearchReq req)(api.post="/api/courses/search"),
  GetCourseDetailResp getCourseDetail(1: GetCourseDetailReq req)(api.get="/api/courses/{course_id}"),
  GetCourseResourceListResp getCourseResourceList(1: GetCourseResourceListReq req)(api.get="/api/courses/{course_id}/resources"),
}


