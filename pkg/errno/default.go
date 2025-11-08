package errno

// 一些常用的的错误, 如果你懒得单独定义也可以直接使用
var (
	Success = NewErrNo(SuccessCode, "成功")

	ParamVerifyError  = NewErrNo(ParamVerifyErrorCode, "参数验证失败")
	ParamMissingError = NewErrNo(ParamMissingErrorCode, "缺少必要参数")

	AuthInvalid             = NewErrNo(AuthInvalidCode, "身份验证失败")
	AuthAccessExpired       = NewErrNo(AuthAccessExpiredCode, "令牌已过期")
	AuthNoToken             = NewErrNo(AuthNoTokenCode, "缺少令牌")
	AuthNoOperatePermission = NewErrNo(AuthNoOperatePermissionCode, "没有操作权限")

	InternalServiceError = NewErrNo(InternalServiceErrorCode, "内部服务错误")
	OSOperationError     = NewErrNo(OSOperateErrorCode, "操作系统调用失败")
	IOOperationError     = NewErrNo(IOOperateErrorCode, "输入输出操作失败")

	QiNiuYunFileError = NewErrNo(QiNiuYunFileErrorCode, "七牛云操作失败")

	//  User Module Errors
	UserPasswordIncorrectError       = NewErrNo(UserPasswordIncorrect, "密码不正确")
	UserPasswordFormatInvalidError   = NewErrNo(UserPasswordFormatInvalid, "密码格式不正确，应为8-20位字母和数字组成")
	UserUsernameFormatInvalidError   = NewErrNo(UserUsernameFormatInvalid, "用户名格式不正确，应为4-16位字母、数字或下划线组成")
	UserEmailFormatInvalidError      = NewErrNo(UserEmailFormatInvalid, "邮箱格式不正确")
	UserVerificationCodeInvalidError = NewErrNo(UserVerificationCodeInvalid, "验证码不正确")
	UserVerificationCodeExpiredError = NewErrNo(UserVerificationCodeExpired, "验证码已过期")
	UserAccountInactiveError         = NewErrNo(UserAccountInactive, "账户未激活")
	UserAccountSuspendedError        = NewErrNo(UserAccountSuspended, "账户已被暂停")

	// Resource Module Errors
	ResourceNotFoundError            = NewErrNo(ResourceNotFound, "资源不存在")
	ResourceAccessDeniedError        = NewErrNo(ResourceAccessDenied, "无权访问该资源")
	ResourceUploadFailedError        = NewErrNo(ResourceUploadFailed, "资源上传失败")
	ResourceDownloadFailedError      = NewErrNo(ResourceDownloadFailed, "资源下载失败")
	ResourceInvalidIDError           = NewErrNo(ResourceInvalidID, "资源ID无效")
	ResourceInvalidRatingError       = NewErrNo(ResourceInvalidRating, "评分必须在0-5之间")
	ResourceInvalidCommentError      = NewErrNo(ResourceInvalidComment, "评论内容不能为空")
	ResourceDuplicateOperationError  = NewErrNo(ResourceDuplicateOperation, "重复操作")
	ResourceReportInvalidReasonError = NewErrNo(ResourceReportInvalidReason, "举报原因不能为空或超过500字符")

	// Course Module Errors
	CourseNotFoundError            = NewErrNo(CourseNotFound, "课程不存在")
	CourseAccessDeniedError        = NewErrNo(CourseAccessDenied, "无权访问该课程")
	CourseInvalidIDError           = NewErrNo(CourseInvalidID, "课程ID无效")
	CourseCommentNotFoundError     = NewErrNo(CourseCommentNotFound, "课程评论不存在")
	CourseRatingNotFoundError      = NewErrNo(CourseRatingNotFound, "课程评分不存在")
	CourseCommentDeleteDeniedError = NewErrNo(CourseCommentDeleteDenied, "无权删除该评论")
	CourseRatingDeleteDeniedError  = NewErrNo(CourseRatingDeleteDenied, "无权删除该评分")

	// Validation Module Errors
	ValidationKeywordTooLongError      = NewErrNo(ValidationKeywordTooLong, "搜索关键词过长")
	ValidationResourceIDInvalidError   = NewErrNo(ValidationResourceIDInvalid, "资源ID无效")
	ValidationCommentTooLongError      = NewErrNo(ValidationCommentTooLong, "评论内容不能超过1000字符")
	ValidationReportReasonTooLongError = NewErrNo(ValidationReportReasonTooLong, "举报原因不能超过500字符")
	ValidationRatingRangeInvalidError  = NewErrNo(ValidationRatingRangeInvalid, "评分必须在0-5之间")
)
