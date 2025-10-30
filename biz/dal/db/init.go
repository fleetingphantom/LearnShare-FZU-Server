package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() error {
	var err error

	DB, err = gorm.Open(mysql.Open("dsn"),
		&gorm.Config{
			PrepareStmt:            true,  // 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
			SkipDefaultTransaction: false, // 不禁用默认事务(即单个创建、更新、删除时使用事务)
			TranslateError:         true,  // 允许翻译错误
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		})
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.InitMySQL mysql connect error: %v", err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.InitMySQL mysql get db error: %v", err))
	}

	sqlDB.SetMaxIdleConns(constants.MaxIdleConns)       // 最大闲置连接数
	sqlDB.SetMaxOpenConns(constants.MaxConnections)     // 最大连接数
	sqlDB.SetConnMaxLifetime(constants.ConnMaxLifetime) // 最大可复用时间
	sqlDB.SetConnMaxIdleTime(constants.ConnMaxIdleTime) // 最长保持空闲状态时间
	DB = DB.WithContext(context.Background())

	if err = sqlDB.Ping(); err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("ping database error: %v", err))
	}

	return nil
}
