package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
)

func CreateUser(ctx context.Context, username, passwordHash, email string) error {
	user := &User{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
		RoleID:       2, // 默认普通用户角色ID
	}

	err := DB.WithContext(ctx).Table(constants.UserTableName).Create(user).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建用户失败: "+err.Error())
	}
	return nil
}

func UpdateUserPassword(ctx context.Context, userID int, newPasswordHash string) error {
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", userID).Update("password_hash", newPasswordHash).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新用户密码失败: "+err.Error())
	}
	return nil
}

func UpdateMajorID(ctx context.Context, userID int, majorID int) error {
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", userID).Update("major_id", majorID).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新用户专业失败: "+err.Error())
	}
	return nil
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询用户失败: "+err.Error())
	}
	return &user, nil
}

func GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", id).First(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询用户失败: "+err.Error())
	}
	return &user, nil
}
