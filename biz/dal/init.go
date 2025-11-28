package dal

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/dal/redis"
)

// Init 初始化数据访问层
func Init() error {
	err := db.Init()
	if err != nil {
		return err
	}

	err = redis.Init()
	if err != nil {
		return err
	}

	return nil
}
