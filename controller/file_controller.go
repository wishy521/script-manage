package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"scripts-manager/common"
)

// ScriptInfo 脚本文件注释信息
type ScriptInfo struct {
	HostLocation  string   `json:"HostLocation"`
	HostIP        []string `json:"HostIP"`
	HostPath      string   `json:"HostPath"`
	HostUser      string   `json:"HostUser"`
	CrontabEnable bool     `json:"CrontabEnable"`
	CrontabData   struct {
		Time    string `json:"Time"`
		Command string `json:"command"`
		Arg     string `json:"arg"`
	} `json:"CrontabData"`
	Language    string `json:"Language"`
	Author      string `json:"Authorer"`
	Description string `json:"Description"`
}

// ExtractCommentBlock 解析并返回注释块中的结构体
func ExtractCommentBlock(content string) (*ScriptInfo, error) {
	// 使用正则表达式匹配注释块
	re := regexp.MustCompile(`(?s): <<COMMENT_BLOCK([\s\S]*?)COMMENT_BLOCK`)
	match := re.FindStringSubmatch(content)
	if len(match) != 2 {
		return nil, fmt.Errorf("no comment block found")
	}

	// 解析JSON到Config结构体
	var info ScriptInfo
	err := json.Unmarshal([]byte(match[1]), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// CreateNewContent 写入剩余内容到新文件
func CreateNewContent(fileContent string, info *ScriptInfo, newFilePath string) error {
	// 重新生成注释
	headerContent := fmt.Sprintf("#!%s\n# Author: %s\n# Description: %s\n", info.Language, info.Author, info.Description)
	var crontabContent string
	if info.CrontabEnable == true {
		crontabContent = fmt.Sprintf("# crontab: %s %s %s",info.CrontabData.Time, info.CrontabData.Command, info.CrontabData.Arg)
	}

	// 移除注释块
	re := regexp.MustCompile(`(?s): <<COMMENT_BLOCK[\s\S]*?COMMENT_BLOCK`)
	scriptContent := re.ReplaceAllString(fileContent, "")


	newContent := headerContent + crontabContent + scriptContent

	// 写入新文件
	err := ioutil.WriteFile(newFilePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write new file: %w", err)
	}

	common.Log.Infof("New file %s created successfully", newFilePath)
	return nil
}


func ExtractContent(filename string) (*[]string, error){
	// 读取原始文件内容
	filePath := fmt.Sprintf("%s%s", common.Conf.Server.Path, filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		common.Log.Errorf("Failed to read file %s: %v", filePath, err)
		return nil, err
	}

	// 提取注释块并解析
	var info *ScriptInfo
	info, err = ExtractCommentBlock(string(content))
	if err != nil {
		common.Log.Errorf("Failed to extract comment block: %v", err)
		return nil, err
	}

	// 写入剩余内容到新文件
	newFilePath := fmt.Sprintf("%snew_%s", common.Conf.Server.Path, filename)
	err = CreateNewContent(string(content), info, newFilePath)
	if err != nil {
		common.Log.Errorf("Failed to write new file: %v", err)
		return nil, err
	}
	return &info.HostIP, nil
}