package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/dal/redis"
	"LearnShare/biz/model/module"
	"LearnShare/biz/model/user"
	oss "LearnShare/pkg"
	"LearnShare/pkg/errno"
	"LearnShare/pkg/utils"
	"context"
	"mime/multipart"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx: ctx, c: c}
}

func (s *UserService) Register(req *user.RegisterReq) error {

	if valid, err := utils.VerifyUsername(req.Username); !valid {
		return err
	}

	if valid, err := utils.VerifyPassword(req.Password); !valid {
		return err
	}

	if valid, err := utils.VerifyEmail(req.Email); !valid {
		return err
	}

	passwordHash, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = passwordHash

	err = db.CreateUser(s.ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) LoginIn(req *user.LoginInReq) (*module.User, error) {
	userInfo, err := db.GetUserByEmail(s.ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := utils.ComparePassword(userInfo.PasswordHash, req.Password); err != nil {
		return nil, err
	}

	userInfo.PasswordHash = ""
	return userInfo.ToUserModule(), nil
}

func (s *UserService) LoginOut() error {
	// 获取当前用户ID用于日志记录

	var errors []error

	uuidStr := GetUuidFormContext(s.c)
	if err := redis.SetBlacklistToken(s.ctx, uuidStr); err != nil {
		errors = append(errors, err)
	}

	// 如果有错误发生，记录日志但不阻止登出
	if len(errors) > 0 {
		// 这里可以添加日志记录
		return errors[0] // 返回第一个错误
	}

	return nil
}

func (s *UserService) SendVerifyEmail(req *user.SendVerifyEmailReq) error {
	code, err := utils.GenerateCode()
	if err != nil {
		return err
	}

	err = utils.MailSendCode(req.Email, code)
	if err != nil {
		return err
	}

	err = redis.PutCodeToCache(s.ctx, req.Email, code)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) VerifyEmail(req *user.VerifyEmailReq) error {

	userId := GetUidFormContext(s.c)

	if !redis.IsKeyExist(s.ctx, req.Email) {
		return errno.UserVerificationCodeExpiredError
	}

	storeCode, err := redis.GetCodeCache(s.ctx, req.Email)
	if err != nil {
		return err
	}

	if storeCode != req.Code {
		return errno.UserVerificationCodeInvalidError
	}

	err = db.UpdateUserStatues(s.ctx, userId, "active")
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateEmail(req *user.UpdateEmailReq) error {

	userId := GetUidFormContext(s.c)

	storeCode, err := redis.GetCodeCache(s.ctx, req.NewEmail)
	if err != nil {
		return err
	}

	if storeCode != req.Code {
		return errno.UserVerificationCodeInvalidError
	}

	err = db.UpdateUserEmail(s.ctx, int(userId), req.NewEmail)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdatePassword(req *user.UpdatePasswordReq) error {
	userId := GetUidFormContext(s.c)
	userInfo, err := db.GetUserByID(s.ctx, userId)
	if err != nil {
		return err
	}

	if err := utils.ComparePassword(userInfo.PasswordHash, req.OldPassword); err != nil {
		return errno.UserPasswordIncorrectError
	}
	newPasswordHash, err := utils.EncryptPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = db.UpdateUserPassword(s.ctx, userInfo.UserID, newPasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateMajor(req *user.UpdateMajorReq) error {
	userId := GetUidFormContext(s.c)
	userInfo, err := db.GetUserByID(s.ctx, userId)
	if err != nil {
		return err
	}

	err = db.UpdateMajorID(s.ctx, int(userInfo.UserID), int(req.NewMajorId))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UploadAvatar(data *multipart.FileHeader) error {
	userId := GetUidFormContext(s.c)

	err := oss.IsImage(data)
	if err != nil {
		return err
	}

	ext := strings.ToLower(path.Ext(data.Filename))

	fileName := strconv.FormatInt(userId, 10) + ext
	storePath := filepath.Join("static", strconv.FormatInt(userId, 10), "avatar")

	if err = oss.SaveFile(data, storePath, fileName); err != nil {
		return err
	}

	url, err := oss.Upload(filepath.Join(storePath, fileName), fileName, "avatar", userId)
	if err != nil {
		return err
	}

	err = db.UpdateAvatarURL(s.ctx, userId, url)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ResetPassword(req *user.ResetPasswordReq) error {
	storeCode, err := redis.GetCodeCache(s.ctx, req.Email)
	if err != nil {
		return err
	}

	if storeCode != req.Code {
		return errno.UserVerificationCodeInvalidError
	}

	newPasswordHash, err := utils.EncryptPassword(req.NewPassword)
	if err != nil {
		return err
	}

	userInfo, err := db.GetUserByEmail(s.ctx, req.Email)
	if err != nil {
		return err
	}

	err = db.UpdateUserPassword(s.ctx, userInfo.UserID, newPasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserInfo(req *user.GetUserInfoReq) (*module.User, error) {
	userInfo, err := db.GetUserByID(s.ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return userInfo.ToUserModule(), nil
}
