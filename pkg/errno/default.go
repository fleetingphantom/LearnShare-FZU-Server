package errno

// 一些常用的的错误, 如果你懒得单独定义也可以直接使用
var (
	Success = NewErrNo(SuccessCode, "ok")

	ParamVerifyError  = NewErrNo(ParamVerifyErrorCode, "parameter validation failed")
	ParamMissingError = NewErrNo(ParamMissingErrorCode, "missing parameter")

	AuthInvalid             = NewErrNo(AuthInvalidCode, "authentication failure")
	AuthAccessExpired       = NewErrNo(AuthAccessExpiredCode, "token expiration")
	AuthNoToken             = NewErrNo(AuthNoTokenCode, "lack of token")
	AuthNoOperatePermission = NewErrNo(AuthNoOperatePermissionCode, "No permission to operate")

	InternalServiceError = NewErrNo(InternalServiceErrorCode, "internal server error")
	OSOperationError     = NewErrNo(OSOperateErrorCode, "os operation failed")
	IOOperationError     = NewErrNo(IOOperateErrorCode, "io operation failed")

	UpYunFileError = NewErrNo(UpYunFileErrorCode, "upyun operation failed")
)
