package service

import (
	"context"
	"testing"
	"time"

	"LearnShare/biz/dal/db"
	redisDal "LearnShare/biz/dal/redis"
	"LearnShare/biz/model/user"
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"LearnShare/pkg/utils"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/cloudwego/hertz/pkg/app"
	goRedis "github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 初始化内存数据库，返回清理函数
func setupTestDB(t *testing.T) func() {
	t.Helper()
	sqliteDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("初始化SQLite失败: %v", err)
	}

	createTableSQL := `
CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    password_hash TEXT,
    email TEXT UNIQUE,
    college_id INTEGER,
    major_id INTEGER,
    avatar_url TEXT,
    reputation_score INTEGER DEFAULT 0,
    role_id INTEGER,
    status TEXT,
    created_at DATETIME,
    updated_at DATETIME
);
`
	if err := sqliteDB.Exec(createTableSQL).Error; err != nil {
		t.Fatalf("创建测试数据表失败: %v", err)
	}

	db.DB = sqliteDB

	return func() {
		sqlDB, err := db.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

// setupTestRedis 初始化内存Redis供单测使用
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, func()) {
	t.Helper()
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("启动 MiniRedis 失败: %v", err)
	}
	redisDal.RDB = goRedis.NewClient(&goRedis.Options{Addr: server.Addr()})

	cleanup := func() {
		_ = redisDal.RDB.Close()
		server.Close()
	}
	return server, cleanup
}

func buildRequestContextWithUser(uid int64) *app.RequestContext {
	ctx := app.NewContext(0)
	ctx.Set(constants.ContextUid, uid)
	return ctx
}

func buildRequestContextWithUserAndUUID(uid int64, uuid string) *app.RequestContext {
	ctx := buildRequestContextWithUser(uid)
	ctx.Set(constants.UUID, uuid)
	return ctx
}

func seedUser(t *testing.T, username, email, password string) db.User {
	t.Helper()
	now := time.Now()
	encrypted, err := utils.EncryptPassword(password)
	if err != nil {
		t.Fatalf("加密密码失败: %v", err)
	}

	userRecord := db.User{
		Username:        username,
		PasswordHash:    encrypted,
		Email:           email,
		ReputationScore: 0,
		RoleID:          2,
		Status:          "inactive",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := db.DB.WithContext(context.Background()).Table(constants.UserTableName).Create(&userRecord).Error; err != nil {
		t.Fatalf("插入用户失败: %v", err)
	}
	return userRecord
}

func TestUserServiceRegister(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()

	svc := NewUserService(context.Background(), nil)
	req := &user.RegisterReq{Username: "user001", Password: "Pass1234", Email: "user001@example.com"}

	if err := svc.Register(req); err != nil {
		t.Fatalf("注册用户失败: %v", err)
	}

	stored, err := db.GetUserByEmail(context.Background(), "user001@example.com")
	if err != nil {
		t.Fatalf("查询注册用户失败: %v", err)
	}
	if utils.ComparePassword(stored.PasswordHash, "Pass1234") != nil {
		t.Fatalf("密码未按预期加密保存")
	}
}

func TestUserServiceLoginIn(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()

	userRecord := seedUser(t, "loginuser", "login@example.com", "Pass1234")

	svc := NewUserService(context.Background(), nil)
	resp, err := svc.LoginIn(&user.LoginInReq{Email: userRecord.Email, Password: "Pass1234"})
	if err != nil {
		t.Fatalf("登录接口返回错误: %v", err)
	}
	if resp.UserId != userRecord.UserID {
		t.Fatalf("预期返回的用户ID为 %d, 实际为 %d", userRecord.UserID, resp.UserId)
	}
}

func TestUserServiceLoginInWrongPassword(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()

	userRecord := seedUser(t, "wrongpass", "wrong@example.com", "Pass1234")

	svc := NewUserService(context.Background(), nil)
	_, err := svc.LoginIn(&user.LoginInReq{Email: userRecord.Email, Password: "Wrong999"})
	if err == nil {
		t.Fatalf("密码错误时应返回错误")
	}
	if errNo, ok := err.(errno.ErrNo); !ok || errNo.ErrorCode != errno.UserPasswordIncorrect {
		t.Fatalf("预期返回密码错误码 %d, 实际错误为 %v", errno.UserPasswordIncorrect, err)
	}
}

func TestUserServiceLoginOut(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()
	_, cleanupRedis := setupTestRedis(t)
	defer cleanupRedis()

	ctx := buildRequestContextWithUserAndUUID(1, "uuid-123")
	svc := NewUserService(context.Background(), ctx)

	if err := svc.LoginOut(); err != nil {
		t.Fatalf("退出登录返回错误: %v", err)
	}

	ok, err := redisDal.IsBlacklistToken(context.Background(), "uuid-123")
	if err != nil {
		t.Fatalf("查询黑名单失败: %v", err)
	}
	if !ok {
		t.Fatalf("刷新黑名单失败")
	}
}

func TestUserServiceVerifyEmail(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()
	server, cleanupRedis := setupTestRedis(t)
	defer cleanupRedis()

	userRecord := seedUser(t, "verify", "verify@example.com", "Pass1234")
	ctx := buildRequestContextWithUser(userRecord.UserID)
	svc := NewUserService(context.Background(), ctx)

	if err := redisDal.PutCodeToCache(context.Background(), userRecord.Email, "654321"); err != nil {
		t.Fatalf("写入验证码失败: %v", err)
	}
	server.FastForward(time.Minute)

	if err := svc.VerifyEmail(&user.VerifyEmailReq{Email: userRecord.Email, Code: "654321"}); err != nil {
		t.Fatalf("校验邮箱失败: %v", err)
	}

	stored, err := db.GetUserByID(context.Background(), userRecord.UserID)
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}
	if stored.Status != "active" {
		t.Fatalf("用户状态应更新为 active, 实际为 %s", stored.Status)
	}
}

func TestUserServiceUpdateEmail(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()
	_, cleanupRedis := setupTestRedis(t)
	defer cleanupRedis()

	userRecord := seedUser(t, "updateEmail", "old@example.com", "Pass1234")
	ctx := buildRequestContextWithUser(userRecord.UserID)
	svc := NewUserService(context.Background(), ctx)

	if err := redisDal.PutCodeToCache(context.Background(), "new@example.com", "111222"); err != nil {
		t.Fatalf("写入验证码失败: %v", err)
	}

	if err := svc.UpdateEmail(&user.UpdateEmailReq{NewEmail: "new@example.com", Code: "111222"}); err != nil {
		t.Fatalf("更新邮箱失败: %v", err)
	}

	stored, err := db.GetUserByID(context.Background(), userRecord.UserID)
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}
	if stored.Email != "new@example.com" {
		t.Fatalf("邮箱未被更新, 当前为 %s", stored.Email)
	}
}

func TestUserServiceUpdatePassword(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()

	userRecord := seedUser(t, "updatePwd", "pwd@example.com", "Pass1234")
	ctx := buildRequestContextWithUser(userRecord.UserID)
	svc := NewUserService(context.Background(), ctx)

	if err := svc.UpdatePassword(&user.UpdatePasswordReq{OldPassword: "Pass1234", NewPassword: "Newpass123"}); err != nil {
		t.Fatalf("更新密码失败: %v", err)
	}

	stored, err := db.GetUserByID(context.Background(), userRecord.UserID)
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}
	if utils.ComparePassword(stored.PasswordHash, "Newpass123") != nil {
		t.Fatalf("新密码未正确保存")
	}
}

func TestUserServiceUpdatePasswordWrongOld(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()

	userRecord := seedUser(t, "wrongOld", "wrongold@example.com", "Pass1234")
	ctx := buildRequestContextWithUser(userRecord.UserID)
	svc := NewUserService(context.Background(), ctx)

	err := svc.UpdatePassword(&user.UpdatePasswordReq{OldPassword: "Wrong999", NewPassword: "Newpass123"})
	if err == nil {
		t.Fatalf("旧密码错误时应返回错误")
	}
	if errNo, ok := err.(errno.ErrNo); !ok || errNo.ErrorCode != errno.UserPasswordIncorrect {
		t.Fatalf("预期错误码 %d, 实际为 %v", errno.UserPasswordIncorrect, err)
	}
}

func TestUserServiceResetPassword(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()
	_, cleanupRedis := setupTestRedis(t)
	defer cleanupRedis()

	userRecord := seedUser(t, "reset", "reset@example.com", "Pass1234")
	if err := redisDal.PutCodeToCache(context.Background(), userRecord.Email, "333444"); err != nil {
		t.Fatalf("写入验证码失败: %v", err)
	}

	svc := NewUserService(context.Background(), nil)
	if err := svc.ResetPassword(&user.ResetPasswordReq{Email: userRecord.Email, Code: "333444", NewPassword: "Reset1234"}); err != nil {
		t.Fatalf("重置密码失败: %v", err)
	}

	stored, err := db.GetUserByEmail(context.Background(), userRecord.Email)
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}
	if utils.ComparePassword(stored.PasswordHash, "Reset1234") != nil {
		t.Fatalf("重置后的密码不匹配")
	}
}

func TestUserServiceGetUserInfo(t *testing.T) {
	cleanupDB := setupTestDB(t)
	defer cleanupDB()

	now := time.Now()
	college := int64(5)
	major := int64(7)
	avatar := "http://example.com/avatar.png"

	userRecord := db.User{
		Username:        "info",
		PasswordHash:    "hash",
		Email:           "info@example.com",
		CollegeID:       &college,
		MajorID:         &major,
		AvatarURL:       &avatar,
		ReputationScore: 10,
		RoleID:          2,
		Status:          "active",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := db.DB.WithContext(context.Background()).Table(constants.UserTableName).Create(&userRecord).Error; err != nil {
		t.Fatalf("插入测试用户失败: %v", err)
	}

	svc := NewUserService(context.Background(), nil)
	result, err := svc.GetUserInfo(&user.GetUserInfoReq{UserID: userRecord.UserID})
	if err != nil {
		t.Fatalf("获取用户信息失败: %v", err)
	}
	if result.UserId != userRecord.UserID {
		t.Fatalf("返回的用户ID不正确")
	}
	if result.AvatarURL != avatar {
		t.Fatalf("头像地址未正确映射")
	}
}
