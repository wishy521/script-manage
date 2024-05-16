package controller

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"scripts-manager/common"
)

// 客户端列表
var clients map[string]ClientInfo

// 初始化客户端列表
func init() {
	clients = make(map[string]ClientInfo)
}

// AddClient 添加客户端到客户端列表
func AddClient(conn net.Conn) {
	clients[conn.RemoteAddr().(*net.TCPAddr).IP.String()] = ClientInfo{
		Addr: conn.RemoteAddr(),
		Conn: conn,
	}
}

// ClientInfo 脚本终端主机连接信息
type ClientInfo struct {
	Addr net.Addr
	Conn net.Conn
}

func TcpConnectionManger(conn net.Conn) {
	defer conn.Close()

	// 获取客户端的地址信息
	clientAddr := conn.RemoteAddr().String()
	common.Log.Infof("Client %s successfully connected", clientAddr)

	// 将客户端添加到客户端列表
	AddClient(conn)

	// 处理客户端发送的消息
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			common.Log.Infof("Client %s disconnected", clientAddr)
			delete(clients, conn.RemoteAddr().(*net.TCPAddr).IP.String())
			return
		}
		common.Log.Infof("Client %s send a message: %s", clientAddr, string(buffer[:n]))
	}
}

func TcpConnectionShow(c *gin.Context)  {
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg": "success",
		"data": clients,
	})
	return
}
