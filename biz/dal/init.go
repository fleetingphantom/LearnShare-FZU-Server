package dal

import "LearnShare/biz/dal/db"

func Init() error {
	err := db.Init()
	if err != nil {
		return err
	}
	return nil
}
