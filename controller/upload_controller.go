package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"scripts-manager/common"
)

func UpLoadFileController(c *gin.Context) {
	// 文件上传
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2001,
			"msg": "failed",
			"data": "上传文件失败",
		})
		common.Log.Error(err.Error())
		return
	}

	// 保存文件
	dst := fmt.Sprintf("%s%s", common.Conf.Server.Path, file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2002,
			"msg": "failed",
			"data": "文件保存失败",
		})
		common.Log.Error(err.Error())
		return
	}
	common.Log.Infof("Successfully received file and saved to %s", dst)

	// 提取文件内容
	var hostArray *[]string
	hostArray, err = ExtractContent(file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2003,
			"msg": "failed",
			"data": "提取文件内容失败",
		})
		common.Log.Errorf("提出文件出错 %s", err.Error())
		return
	}
	newFilePath := fmt.Sprintf("%snew_%s", common.Conf.Server.Path, file.Filename)
	var hostIPMap = make(map[string]bool)
	var msgData = ""
	hostIPMap, msgData, err = SendFileToHost(hostArray, newFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2004,
			"msg": msgData,
			"data": "发送文件失败",
		})
		common.Log.Errorf("发送文件出错 %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg": "success",
		"data": hostIPMap,
	})
	common.Log.Info("发送文件任务完成")
	return
}