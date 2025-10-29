namespace go module

struct BaseResp{
    required i32 code,
    required string message,
}

struct User{
    required i64 userId,
    required string userName,
    optional string password,
    required string email,
    required i64 collegeId,
    required i64 majorId,
    required string avatarUrl,
    required i64 reputationScore,
    required i64 roleId,
    required string status,
    required i64 createdAt,
    required i64 updatedAt,
}

struct Course{
    required i64 courseId,
    required string courseName,
    required i64 teacherId,
    required i64 credit,
    required i64 majorId,
    required string grade,
    required string description,
    required i64 createdAt,
    required i64 updatedAt,
}

struct CourseRating{
    required i64 ratingId,
    required i64 userId,
    required i64 courseId,
    required i64 recommendation,
    required string difficulty,
    required i64 workload,
    required i64 usefulness,
    required bool isVisible,
    required i64 createdAt,
}


struct CourseComment{
    required i64 commentId,
    required i64 userId,
    required i64 courseId,
    required string content,
    required i64 parentId,
    required i64 likes,
    required bool isVisible,
    required i64 status,
    required i64 createdAt,
}


struct Resource{
    required i64 resourceId,
    required string resourceName,
    required string description,
    required string resourceUrl,
    required string type,
    required i64 size,
    required i64 uploaderId,
    required i64 courseId,
    required i64 downloadCount,
    required i64 averageRating,
    required i64 ratingCount,
    required string status,
    required i64 createdAt,
}

struct ResourceRating{
    required i64 ratingId,
    required i64 userId,
    required i64 resourceId,
    required i64 recommendation,
    required bool isVisible,
    required i64 createdAt,
}


struct ResourceComment{
    required i64 commentId,
    required i64 userId,
    required i64 resourceId,
    required string content,
    required i64 parentId,
    required i64 likes,
    required bool isVisible,
    required i64 status,
    required i64 createdAt,
}

