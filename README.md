# 脚本管理系统服务端

## 简介
这是一个linux主机脚本文件管理系统；本质上就是实现了服务端向客户端发送文件的一个工具。
正常修改Linux脚本内容，需要使用远程工具连接服务器，然后用Linux自带的编辑工具修改脚本有时候效率会比较低，因此脚本一般都是在工作电脑上完成后上传至服务器。
在工作电脑上编辑完脚本后，仍需要打开远程连接工具，然后上传到指定路径。



## 使用方式

### 1.使用一台Linux主机部署script-manage服务作为服务端
```bash
git clone https://github.com/wishy521/script-manage.git
cd script-manage
# 构建项目
go mod tidy
go build -o scriptmanage main.go
# 启动
nohup ./scriptmanage 2>&1 &
## 注意服务端会监听两个端口，一个专门用于Linux主机连接，另一负责应用相关的，可以在config.yml中修改端口号
```

### 2.需要脚本被管理的主机(Linux)运行script-client客户端工具
```bash
git clone https://github.com/wishy521/script-client.git
cd script-client
# 构建项目
go mod tidy
go build -o scriptclient main.go
# 启动（注意端口号为服务器端server.port.tcp的端口号）
nohup ./scriptclient --address=[服务端IP:端口] > scripts-client.log 2>&1 &
```

### 3.在原有的脚本中添加注释块(参考模板example_shell.sh)
注释块是一个json格式的，方便解析里面的内容
```json
{
    "HostLocation" : "HuaweiCloud",
    "HostIP" : ["10.10.10.10","10.10.10.11","127.0.0.1"],
    "FileInfo" : {
        "Path" : "shell.sh",
        "Owner" : "root",
        "Group" : "root",
        "Perm" : "0400"
    },
    "CrontabEnable" : true,
    "CrontabData" : {
        "Time" : "0 1 * * *",
        "command" : "/usr/bin/sh",
        "arg" : "/home/shell/shell.sh"
    },
    "Language" : "bash",
    "Author" : "luojiashuo",
    "Description" : "shell脚本注释模板"
}
```
### 4.通过服务端的文件上传接口上传脚本文件就会被转发到指定Linux主机

```bash
curl -XPOST --user script:script123 -H "Content-Type:multipart/form-data"  -F "file=@shell.sh" http://127.0.0.1:7080/upload  
```
参数解析
```
curl -XPOST                            # 使用post请求 
--user script:script123                # user为服务端的auth账号密码
-H "Content-Type:multipart/form-data" 
-F "file=@shell.sh"                    # 指定上传文件
http://127.0.0.1:7080/upload           # 指定服务端地址
```
相应结果说明
```
{
    "code": 2000,  # 表示请求成功
    "data": {
        "10.10.10.10": false, # 主机不存在
        "10.10.10.11": true,  # 成功转发的主机
        "127.0.0.1": false
    },
    "msg": "warn"
}
```

### 5.使用gitlab+gitlab-runner管理脚本(扩展)
```
scripts-prod:
  stage: scripts
  tags:
  - prod
  script:
    - |
      changed_files=$(git diff --name-only HEAD~1 HEAD)
      for file in $changed_files; do
        if [[ $file == *.sh ]]; then
          if [ -f $file ]; then
            curl -XPOST --user script:script123 -H "Content-Type:multipart/form-data"  -F "file=$file" http://127.0.0.1:7080/upload
          fi
        fi
      done
  when: manual
  only:
    refs:
      - master
    changes:
      - scripts/*
```