package user

import (
	"Tally/common"
	"Tally/config"
	"Tally/global"
	"archive/zip"
	"github.com/labstack/echo/v4"
	"io"
	"os"
	"path"
)

// ExportLog 压缩日志
func ExportLog(c echo.Context) error {
	//进行压缩
	//读取文件夹
	dir, err := os.ReadDir(config.Config.Logs.Path)
	if err != nil {
		return common.Fail(c, global.UserCode, global.CreateLogErr)
	}
	// 创建zip文件
	zipFile, err := os.Create(config.Config.Logs.Path + "logfile.zip")
	if err != nil {
		global.Global.Log.Error(err)
		return common.Fail(c, global.LogCode, global.CreateLogErr)
	}
	zips := zip.NewWriter(zipFile)
	for _, res := range dir {
		if path.Ext(res.Name()) == ".zip" {
			continue
		}
		file, err := os.Open(config.Config.Logs.Path + res.Name())
		if err != nil {
			global.Global.Log.Error(err)
			return common.Fail(c, global.LogCode, global.CreateLogErr)
		}
		create, err := zips.Create(res.Name())

		if err != nil {
			global.Global.Log.Error(err)
			return common.Fail(c, global.LogCode, global.CreateLogErr)
		}
		_, err = io.Copy(create, file)
		if err != nil {
			global.Global.Log.Error(err)
			return common.Fail(c, global.LogCode, global.CreateLogErr)
		}
		err = file.Close()
		if err != nil {
			global.Global.Log.Error(err)
			err := zipFile.Close()
			if err != nil {
				return err
			}
			return common.Fail(c, global.LogCode, global.CreateLogErr)
		}
	}
	err = zipFile.Close()
	if err != nil {
		return err
	}
	return common.Picture(c, config.Config.Logs.Path+"logfile.zip")
}
