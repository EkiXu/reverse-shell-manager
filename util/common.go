package util

import (
	"sh.ieki.xyz/global"
	"os"
)

// @title    PathExists
// @description   文件目录是否存在
// @auth
// @param     path            string
// @return    err             error

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// @title    createDir
// @description   批量创建文件夹
// @auth
// @param     dirs            string
// @return    err             error

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.SERVER_LOG.Debug("create directory ", v)
			err = os.MkdirAll(v, os.ModePerm)
			if err != nil {
				global.SERVER_LOG.Error("create directory", v, " error:", err)
			}
		}
	}
	return err
}
