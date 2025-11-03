package utils

import (
	"LearnShare/config"
	"errors"
	"strconv"
	"strings"
)

func GetMysqlDSN() (string, error) {
	if config.Mysql == nil {
		return "", errors.New("未找到配置")
	}

	dsn := strings.Join([]string{
		config.Mysql.Username, ":", config.Mysql.Password,
		"@tcp(", config.Mysql.Addr, ")/",
		config.Mysql.Database, "?charset=" + config.Mysql.Charset + "&parseTime=true",
	}, "")

	return dsn, nil
}

func GetServerAddress() string {
	if config.Server == nil {
		panic("未找到配置")
		return ""
	}

	address := strings.Join([]string{
		config.Server.Addr, ":", strconv.Itoa(config.Server.Port),
	}, "")

	return address
}
