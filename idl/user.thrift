namespace go user
include "model.thrift"

//用户注册
struct RegisterReq {
  required string username;
  required string password;
  required string email;
}
struct RegisterResp {
  required model.BaseResp baseResponse;
}

//用户登录
struct LoginInReq {
  required string email;
  required string password;
}
struct LoginInResp {
  required model.BaseResp baseResponse;
  optional model.User user;
}

//用户登出
struct LoginOutReq {
}
struct LoginOutResp {
  required model.BaseResp baseResponse;
}

//获取邮箱验证码
struct SendVerifyEmailReq {
  required string email;
}
struct SendVerifyEmailResp {
  required model.BaseResp baseResponse;
}

//验证邮箱验证码
struct VerifyEmailReq {
  required string email;
  required string code;
}
struct VerifyEmailResp {
  required model.BaseResp baseResponse;
}

//修改邮箱
struct updateEmailReq {
  required string new_email;
  required string code;
}
struct updateEmailResp {
  required model.BaseResp baseResponse;
}

//修改密码
struct UpdatePasswordReq {
  required string old_password;
  required string new_password;
}
struct UpdatePasswordResp {
  required model.BaseResp baseResponse;
}

//修改专业
struct updateMajorReq {
  required i64 new_majorId;
}
struct updateMajorResp {
  required model.BaseResp baseResponse;
}


//上传头像
struct uploadAvatarReq {
//  required binary avatar;
}
struct uploadAvatarResp {
  required model.BaseResp baseResponse;
}


//重置密码
struct ResetPasswordReq {
  required string email;
  required string newPassword;
  required string code;
}
struct ResetPasswordResp {
  required model.BaseResp baseResponse;
}


//刷新Token
struct RefreshTokenReq {
}
struct RefreshTokenResp {
  required model.BaseResp baseResponse;
}

//获取用户信息
struct GetUserInfoReq {
  required i64 user_id (api.path="user_id");
}
struct GetUserInfoResp {
  required model.BaseResp baseResponse;
  optional model.User user;
}


service UserService {
  RegisterResp register(1: RegisterReq req)(api.post="/api/auth/register"),
  LoginInResp loginIn(1: LoginInReq req)(api.post="/api/auth/login"),
  LoginOutResp loginOut(1: LoginOutReq req)(api.post="/api/auth/logout"),
  SendVerifyEmailResp sendVerifyEmail(1: SendVerifyEmailReq req)(api.post="/api/users/me/email/get"),
  VerifyEmailResp verifyEmail(1: VerifyEmailReq req)(api.post="/api/users/me/email/verify"),
  updateEmailResp updateEmail(1: updateEmailReq req)(api.put="/api/users/me/email"),
  UpdatePasswordResp updatePassword(1: UpdatePasswordReq req)(api.put="/api/users/me/password"),
  updateMajorResp updateMajor(1: updateMajorReq req)(api.put="/api/users/me/major"),
  uploadAvatarResp uploadAvatar(1: uploadAvatarReq req)(api.put="/api/users/avatar"),
  ResetPasswordResp resetPassword(1: ResetPasswordReq req)(api.post="/api/users/me/password/reset"),
  RefreshTokenResp refreshToken(1: RefreshTokenReq req)(api.post="/api/auth/refresh"),
  GetUserInfoResp getUserInfo(1: GetUserInfoReq req)(api.get="/api/users/:user_id"),
}

