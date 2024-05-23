package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"scripts-manage/common"
)

// ScriptInfo 脚本文件注释信息
type ScriptInfo struct {
	HostLocation string   `json:"HostLocation"`
	HostIP       []string `json:"HostIP"`
	FileInfo     struct {
		Path  string `json:"Path"`
		Owner string `json:"Owner"`
		Group string `json:"Group"`
		Perm  string `json:"Perm"`
	} `json:"FileInfo"`
	CrontabEnable bool `json:"CrontabEnable"`
	CrontabData   struct {
		Time    string `json:"Time"`
		Command string `json:"command"`
		Arg     string `json:"arg"`
	} `json:"CrontabData"`
	Language    string `json:"Language"`
	Author      string `json:"Author"`
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
		common.Log.Errorf(err.Error())
		return nil, err
	}
	return &info, nil
}

func ExtractContent(filePath string) (*[]string, error) {
	// 读取原始文件内容
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

	return &info.HostIP, nil
}
