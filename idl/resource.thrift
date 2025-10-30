namespace go resource
include "model.thrift"

// 搜索资源请求
struct SearchResourceReq {
    1: optional string keyword,          
    2: optional i64 tagId,               
    3: optional string sortBy, 
    4: optional i64 courseId,  
    5: required i32 pageSize, 
    6: required i32 pageNum,             
}

struct SearchResourceResp {
    1: required model.BaseResp baseResp,
    2: required list<model.Resource> resources,
    3: required i32 total, 
}

// 上传资源请求
struct UploadResourceReq {
    1: required binary fileData, 
    2: required string title,  
    3: optional string description,
    4: required i64 courseId, 
    5: optional list<string> tags,
}

struct UploadResourceResp {
    1: required model.BaseResp baseResp,
    2: optional model.Resource resource,
}

// 下载资源请求
struct DownloadResourceReq {
    1: required i64 resourceId, 
}

struct DownloadResourceResp {
    1: required model.BaseResp baseResp,
    2: required string downloadUrl,
}

// 举报资源请求
struct ReportResourceReq {
    1: required i64 resourceId,
    2: required string reason, 
}

struct ReportResourceResp {
    1: required model.BaseResp baseResp,
}

// 获取资源信息请求
struct GetResourceReq {
    1: required i64 resourceId, 
}

struct GetResourceResp {
    1: required model.BaseResp baseResp,
    2: optional model.Resource resource,
}

// 资源服务
service ResourceService {
    SearchResourceResp searchResources(1: SearchResourceReq req)(api.get="/api/resources/search"),
    UploadResourceResp uploadResource(1: UploadResourceReq req)(api.post="/api/resources"),
    DownloadResourceResp downloadResource(1: DownloadResourceReq req)(api.get="/api/resources/{resource_id}/download"),
    ReportResourceResp reportResource(1: ReportResourceReq req)(api.post="/api/resources/{resource_id}/report"),
    GetResourceResp getResource(1: GetResourceReq req)(api.get="/api/resources/{resource_id}"),
}