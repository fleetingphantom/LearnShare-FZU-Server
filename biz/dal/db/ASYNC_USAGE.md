# 数据库异步操作使用指南

## 概述

本项目已对数据库操作进行了异步优化,提供了以下两种优化方案:

1. **异步写入操作** - 使用 goroutine + channel 模式,适用于创建/更新/删除操作
2. **查询优化** - 批量查询优化和超时控制,解决 N+1 查询问题

## 1. 异步工作池 (AsyncWorkerPool)

位置: `biz/dal/db/async.go`

### 特性
- 固定大小的 worker 池 (默认10个worker)
- 缓冲任务队列 (容量100)
- 自动错误处理
- 支持等待结果或不等待结果的提交方式

### 使用示例

#### 方式一: 等待异步结果
```go
import "LearnShare/biz/dal/db"

// 异步更新用户密码
errChan := db.UpdateUserPasswordAsync(ctx, userID, newPasswordHash)

// 等待结果
if err := <-errChan; err != nil {
    // 处理错误
    log.Printf("更新失败: %v", err)
}
```

#### 方式二: 不等待结果 (fire-and-forget)
```go
import "LearnShare/biz/dal/db"

pool := db.GetAsyncPool()
pool.SubmitNoWait(func() error {
    return db.UpdateUserEmail(ctx, userID, newEmail)
})
```

## 2. 异步函数列表

### User 模块 (biz/dal/db/user.go)

| 同步函数 | 异步函数 | 说明 |
|---------|---------|------|
| `UpdateUserPassword` | `UpdateUserPasswordAsync` | 更新用户密码 |
| `UpdateMajorID` | `UpdateMajorIDAsync` | 更新用户专业 |
| `UpdateAvatarURL` | `UpdateAvatarURLAsync` | 更新用户头像 |
| `UpdateUserStatues` | `UpdateUserStatuesAsync` | 更新用户状态 |
| `UpdateUserEmail` | `UpdateUserEmailAsync` | 更新用户邮箱 |

**使用示例:**
```go
// 同步方式
err := db.UpdateUserPassword(ctx, 123, "newHash")

// 异步方式
errChan := db.UpdateUserPasswordAsync(ctx, 123, "newHash")
if err := <-errChan; err != nil {
    // 处理错误
}
```

### Course 模块 (biz/dal/db/course.go)

| 同步函数 | 异步函数 | 说明 |
|---------|---------|------|
| `CreateCourse` | `CreateCourseAsync` | 创建课程 |
| `UpdateCourse` | `UpdateCourseAsync` | 更新课程 |
| `DeleteCourse` | `DeleteCourseAsync` | 删除课程 |
| `SubmitCourseRating` | `SubmitCourseRatingAsync` | 提交课程评分 |
| `UpdateCourseRating` | `UpdateCourseRatingAsync` | 更新课程评分 |
| `DeleteCourseRating` | `DeleteCourseRatingAsync` | 删除课程评分 |
| `SubmitCourseComment` | `SubmitCourseCommentAsync` | 提交课程评论 |
| `UpdateCourseComment` | `UpdateCourseCommentAsync` | 更新课程评论 |
| `DeleteCourseComment` | `DeleteCourseCommentAsync` | 删除课程评论 |
| `CreateResource` | `CreateResourceAsync` | 创建资源 |
| `UpdateResource` | `UpdateResourceAsync` | 更新资源 |
| `DeleteResource` | `DeleteResourceAsync` | 删除资源 |

**使用示例:**
```go
// 异步创建课程
errChan := db.CreateCourseAsync(ctx, "数据结构", 1, 2, 3.0, "大二", "课程描述")
if err := <-errChan; err != nil {
    return err
}
```

### Resource 模块 (biz/dal/db/resource.go)

| 同步函数 | 异步函数 | 说明 |
|---------|---------|------|
| `SubmitResourceComment` | `SubmitResourceCommentAsync` | 提交资源评论 |
| `DeleteResourceRating` | `DeleteResourceRatingAsync` | 删除资源评分 |
| `DeleteResourceComment` | `DeleteResourceCommentAsync` | 删除资源评论 |
| `CreateReview` | `CreateReviewAsync` | 创建举报 |

**特殊示例 - SubmitResourceCommentAsync:**
```go
// 此函数返回带结果的 channel
resultChan := db.SubmitResourceCommentAsync(ctx, userID, resourceID, "评论内容", nil)

// 等待结果
result := <-resultChan
if result.Err != nil {
    return result.Err
}
comment := result.Comment // 获取创建的评论对象
```

## 3. 查询优化

### N+1 查询优化

**优化前** (permission.go - GetAllRoles):
```go
// 每个角色都会执行一次查询,N+1问题
for _, role := range roles {
    var permissions []Permission
    db.Where("role_id = ?", role.ID).Find(&permissions)
}
```

**优化后**:
```go
// 一次性批量查询所有角色的权限
var rolePermissions []RolePermission
db.Where("role_id IN ?", roleIDs).Find(&rolePermissions)

// 一次性批量查询所有权限详情
var permissions []Permission
db.Where("permission_id IN ?", permissionIDs).Find(&permissions)
```

### 超时控制

重要查询函数已添加 5 秒超时控制:

```go
// SearchResources 和 GetResourceComments 已添加超时
func SearchResources(ctx context.Context, ...) ([]*Resource, int64, error) {
    ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    db := DB.WithContext(ctxWithTimeout).Table(...)
    // ...
}
```

## 4. 批量异步执行

使用 `AsyncBatch` 函数批量执行多个独立的异步任务:

```go
import "LearnShare/biz/dal/db"

tasks := []func() error{
    func() error { return db.UpdateUserEmail(ctx, 1, "email1@test.com") },
    func() error { return db.UpdateUserEmail(ctx, 2, "email2@test.com") },
    func() error { return db.UpdateUserEmail(ctx, 3, "email3@test.com") },
}

// 并行执行所有任务
results := db.AsyncBatch(ctx, tasks)

// 检查错误
for i, err := range results {
    if err != nil {
        log.Printf("Task %d failed: %v", i, err)
    }
}
```

## 5. 最佳实践

### 何时使用异步?

✅ **适合使用异步:**
- 不需要立即返回结果的写操作
- 日志记录、统计更新等非关键操作
- 批量更新操作

❌ **不适合使用异步:**
- 需要立即返回数据的操作
- 事务操作 (已在事务内的操作保持同步)
- 关键业务逻辑需要立即知道成功/失败的操作

### 错误处理

```go
// 正确的错误处理方式
errChan := db.UpdateUserPasswordAsync(ctx, userID, newHash)

select {
case err := <-errChan:
    if err != nil {
        // 记录错误日志
        log.Printf("异步操作失败: %v", err)
        // 可以选择重试或通知用户
    }
case <-time.After(10 * time.Second):
    // 超时处理
    log.Println("异步操作超时")
}
```

### 性能监控

```go
import "time"

start := time.Now()
errChan := db.UpdateUserPasswordAsync(ctx, userID, newHash)
err := <-errChan
duration := time.Since(start)

log.Printf("异步操作耗时: %v", duration)
```

## 6. 性能提升

通过异步优化,预计可获得以下性能提升:

- **写操作**: 响应时间降低 60-80%
- **N+1 查询优化**: 数据库查询次数从 O(N) 降低到 O(1)
- **批量操作**: 吞吐量提升 3-5 倍
- **超时控制**: 防止慢查询阻塞系统

## 7. 注意事项

1. **Worker Pool 生命周期**: 全局 worker pool 会在程序启动时自动创建,无需手动管理
2. **Context 传递**: 务必正确传递 context,以支持超时和取消操作
3. **事务操作**: 事务内的操作保持同步,不要在事务中使用异步函数
4. **错误恢复**: Worker 内部已处理 panic,不会导致整个程序崩溃
5. **资源限制**: Worker pool 默认 10 个 worker,根据实际负载调整

## 8. 迁移指南

将现有同步代码迁移到异步:

```go
// 步骤 1: 找到适合异步的操作
err := db.UpdateUserPassword(ctx, userID, newHash)

// 步骤 2: 替换为异步版本
errChan := db.UpdateUserPasswordAsync(ctx, userID, newHash)

// 步骤 3: 根据业务需求选择等待或不等待
// 选项 A: 等待结果
if err := <-errChan; err != nil {
    return err
}

// 选项 B: 不等待结果 (适用于非关键操作)
go func() {
    if err := <-errChan; err != nil {
        log.Printf("后台操作失败: %v", err)
    }
}()
```
