package main

import (
	"fmt"
	"net"
	"net/http"
	"scripts-manage/common"
	"scripts-manage/controller"
)

func main() {
	// 加载配置文件到全局配置结构体
	common.InitConfig()

	// 初始化日志
	common.InitLogger()

	// 监听端口
	tcpPort := common.Conf.Server.Port.Tcp
	httpPort := common.Conf.Server.Port.Http
	router := controller.InitRoutes()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: router,
	}

	// 在goroutine中初始化服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Log.Errorf("listen: %s\n", err)
		}
	}()

	common.Log.Infof(fmt.Sprintf("Web Server is running at %d", httpPort))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", tcpPort))
	if err != nil {
		common.Log.Errorf("TCP server listen failed : %s", err)
		return
	}
	defer listener.Close()

	common.Log.Infof(fmt.Sprintf("TCP Server is running at %d", tcpPort))

	// 循环接受客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			common.Log.Errorf("Receive client connection failed: %s", err)
			continue
		}
		go controller.HandleConnection(conn)
	}
}
