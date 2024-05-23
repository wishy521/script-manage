package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"scripts-manage/common"
)

func UpLoadFileController(c *gin.Context) {
	// 文件上传
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2001,
			"msg":  "failed",
			"data": "上传文件失败",
		})
		common.Log.Errorf("UpLoad file failed %s", err)
		return
	}

	// 保存文件
	filePath := fmt.Sprintf("%s%s", common.Conf.Server.Path, file.Filename)
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2002,
			"msg":  "failed",
			"data": "文件保存失败",
		})
		common.Log.Errorf("Save uploaded file failed %s", err)
		return
	}
	common.Log.Infof("Successfully received file and saved to %s", filePath)

	// 提取发送地址
	hostArray, err := ExtractContent(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2003,
			"msg":  "failed",
			"data": "提取文件内容失败",
		})
		common.Log.Errorf("Extract file content failed %s", err)
		return
	}

	hostIPMap, msgData, err := SendFileToHost(hostArray, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2004,
			"msg":  "failed",
			"data": "发送文件失败",
		})
		common.Log.Errorf("Sending file failed %s", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  msgData,
		"data": hostIPMap,
	})
	common.Log.Info("Sending file task complete")
	return
}
