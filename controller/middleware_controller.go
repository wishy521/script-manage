package controller

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"scripts-manage/common"
	"strings"
	"sync"
	"time"
)

// WhiteListMiddleware 白名单认证
func WhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		allowed := false

		// 遍历白名单
		whitelist := common.Conf.Server.Whitelist
		for _, cidr := range whitelist {
			// 解析IP或CIDR
			if strings.Contains(cidr, "/") {
				_, subnet, err := net.ParseCIDR(cidr)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"code": 1001,
						"msg":  "failed",
						"data": "Invalid CIDR",
					})
					common.Log.Errorf("Invalid CIDR %s", err)
					return
				}
				if subnet.Contains(net.ParseIP(clientIP)) {
					allowed = true
					break
				}
			} else {
				// 处理单个IP的情况
				if clientIP == cidr {
					allowed = true
					break
				}
			}
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 1002,
				"msg":  "failed",
				"data": "address is invalid",
			})
			return
		}
		c.Next()
	}
}

// BasicAuthMiddleware 基础认证中间件
func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 1101,
				"msg":  "failed",
				"data": "Authorization header is missing",
			})
			return
		}

		// 解析认证信息
		if !strings.HasPrefix(authHeader, "Basic ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header",
			})
			return
		}

		// 解码 Base64 编码的认证信息
		decoded, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  1102,
				"msg":   "failed",
				"error": "Invalid authorization header",
			})
			return
		}

		// 分割用户名和密码
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  1103,
				"msg":   "failed",
				"error": "Invalid authorization header",
			})
			return
		}

		// 验证用户名和密码
		username := common.Conf.Server.Auth.Username
		password := common.Conf.Server.Auth.Password
		if parts[0] != username || parts[1] != password {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  1104,
				"msg":   "failed",
				"error": "Invalid username or password",
			})
			return
		}
		c.Next()
	}
}

// RateLimiterMiddleware 限流中间件
func RateLimiterMiddleware() gin.HandlerFunc {
	var mu sync.Mutex
	ipRequests := make(map[string][]time.Time)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		currentTime := time.Now()

		mu.Lock()
		requestTimes, exists := ipRequests[clientIP]
		duration := time.Duration(common.Conf.Server.Limit.Duration) * time.Minute

		if exists {
			// 移除过期的请求时间
			var validTimes []time.Time
			for _, t := range requestTimes {
				if currentTime.Sub(t) < duration {
					validTimes = append(validTimes, t)
				}
			}
			ipRequests[clientIP] = validTimes

			// 检查请求次数是否超过限制
			if len(validTimes) >= common.Conf.Server.Limit.Count {
				mu.Unlock()
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"code": 1201,
					"msg":  "failed",
					"data": "Too many requests",
				})
				return
			}
			ipRequests[clientIP] = append(ipRequests[clientIP], currentTime)
		} else {
			// 新IP地址，初始化请求时间
			ipRequests[clientIP] = []time.Time{currentTime}
		}
		mu.Unlock()
		c.Next()
	}
}
