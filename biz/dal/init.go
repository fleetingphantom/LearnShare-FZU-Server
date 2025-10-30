package dal

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/dal/redis"
)

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
