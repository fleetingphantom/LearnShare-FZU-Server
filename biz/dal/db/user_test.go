package db

import (
	"context"
	"testing"
	"time"

	"LearnShare/pkg/constants"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// initTestDB 用于初始化内存数据库，返回清理函数
func initTestDB(t *testing.T) func() {
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

	DB = sqliteDB

	return func() {
		sqlDB, err := DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

func insertUser(t *testing.T, username, email, passwordHash string) User {
	t.Helper()
	now := time.Now()
	user := User{
		Username:        username,
		PasswordHash:    passwordHash,
		Email:           email,
		ReputationScore: 0,
		RoleID:          2,
		Status:          "inactive",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := DB.WithContext(context.Background()).Table(constants.UserTableName).Create(&user).Error; err != nil {
		t.Fatalf("插入测试用户失败: %v", err)
	}
	return user
}

func TestCreateUser(t *testing.T) {
	cleanup := initTestDB(t)
	defer cleanup()

	ctx := context.Background()
	err := CreateUser(ctx, "user1234", "hashpwd", "user@example.com")
	if err != nil {
		t.Fatalf("创建用户返回错误: %v", err)
	}

	var u User
	if err := DB.WithContext(ctx).Table(constants.UserTableName).Where("email = ?", "user@example.com").First(&u).Error; err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}

	if u.Username != "user1234" {
		t.Fatalf("预期用户名为 user1234, 实际为 %s", u.Username)
	}
	if u.RoleID != 2 {
		t.Fatalf("预期角色ID为2, 实际为 %d", u.RoleID)
	}
	if u.Status != "inactive" {
		t.Fatalf("预期状态为 inactive, 实际为 %s", u.Status)
	}
}

func TestUpdateUserPassword(t *testing.T) {
	cleanup := initTestDB(t)
	defer cleanup()

	inserted := insertUser(t, "tester", "tester@example.com", "old")

	ctx := context.Background()
	err := UpdateUserPassword(ctx, inserted.UserID, "newhash")
	if err != nil {
		t.Fatalf("更新密码失败: %v", err)
	}

	var u User
	if err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", inserted.UserID).First(&u).Error; err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}
	if u.PasswordHash != "newhash" {
		t.Fatalf("预期密码哈希为 newhash, 实际为 %s", u.PasswordHash)
	}
}

func TestGetUserByEmail(t *testing.T) {
	cleanup := initTestDB(t)
	defer cleanup()

	insertUser(t, "alice", "alice@example.com", "hash")

	ctx := context.Background()
	user, err := GetUserByEmail(ctx, "alice@example.com")
	if err != nil {
		t.Fatalf("根据邮箱获取用户失败: %v", err)
	}
	if user.Username != "alice" {
		t.Fatalf("预期用户名为 alice, 实际为 %s", user.Username)
	}
}

func TestUpdateUserEmailAndStatus(t *testing.T) {
	cleanup := initTestDB(t)
	defer cleanup()

	inserted := insertUser(t, "bob", "bob@example.com", "hash")
	ctx := context.Background()

	if err := UpdateUserEmail(ctx, int(inserted.UserID), "new@example.com"); err != nil {
		t.Fatalf("更新邮箱失败: %v", err)
	}

	if err := UpdateUserStatues(ctx, inserted.UserID, "active"); err != nil {
		t.Fatalf("更新状态失败: %v", err)
	}

	var u User
	if err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", inserted.UserID).First(&u).Error; err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}
	if u.Email != "new@example.com" {
		t.Fatalf("预期邮箱为 new@example.com, 实际为 %s", u.Email)
	}
	if u.Status != "active" {
		t.Fatalf("预期状态为 active, 实际为 %s", u.Status)
	}
}

func TestUpdateUserAvatarAndMajor(t *testing.T) {
	cleanup := initTestDB(t)
	defer cleanup()

	inserted := insertUser(t, "carol", "carol@example.com", "hash")
	ctx := context.Background()

	if err := UpdateAvatarURL(ctx, inserted.UserID, "http://example.com/avatar.png"); err != nil {
		t.Fatalf("更新头像失败: %v", err)
	}
	if err := UpdateMajorID(ctx, int(inserted.UserID), 3); err != nil {
		t.Fatalf("更新专业失败: %v", err)
	}

	var u User
	if err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", inserted.UserID).First(&u).Error; err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}
	if u.AvatarURL == nil || *u.AvatarURL != "http://example.com/avatar.png" {
		t.Fatalf("头像地址未正确更新")
	}
	if u.MajorID == nil || *u.MajorID != 3 {
		t.Fatalf("专业ID未正确更新")
	}
}
