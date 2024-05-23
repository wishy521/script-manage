package controller

import (
	"net"
	"os"
	"scripts-manage/common"
	"time"
)

func SendFileToHost(hostArray *[]string, filePath string) (map[string]bool, string, error) {
	hostIPMap := make(map[string]bool)
	successCount := 0
	sendInfo := ""
	for _, host := range *hostArray {
		client, ok := clients[host]
		if ok {
			err := WriteConnection(client.Conn, filePath)
			if err != nil {
				common.Log.Errorf("write to host %s connection failed ", host)
				return nil, sendInfo, err
			}
			hostIPMap[host] = true
			successCount++
		} else {
			common.Log.Warnf("host %s connection dose not exeit ", host)
			hostIPMap[host] = false
		}
	}
	switch {
	case successCount == len(*hostArray):
		sendInfo = "success"
	case 0 < successCount && successCount < len(*hostArray):
		sendInfo = "warn"
	case successCount == 0:
		sendInfo = "failed"
	default:
		sendInfo = "error"
	}
	return hostIPMap, sendInfo, nil
}

// WriteConnection 写入内容到TCP连接
func WriteConnection(conn net.Conn, filePath string) error {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		common.Log.Errorf("Error opening file: %s", err.Error())
		return err
	}
	defer file.Close()

	// 发送 START 标识符
	_, err = conn.Write([]byte("START"))
	if err != nil {
		common.Log.Errorf("Error sending START: %s", err.Error())
		return err
	}
	time.Sleep(100 * time.Millisecond)
	// 从文件中读取内容并写入连接
	buffer := make([]byte, 1024)
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			break
		}
		conn.Write(buffer[:bytesRead])
	}
	time.Sleep(100 * time.Millisecond)
	// 发送 END 标识符
	_, err = conn.Write([]byte("END"))
	if err != nil {
		common.Log.Errorf("Error sending END: %s", err.Error())
		return err
	}
	return nil
}
