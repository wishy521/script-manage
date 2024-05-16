: <<COMMENT_BLOCK
{
    "HostLocation" : "HuaweiCloud",
    "HostIP" : ["10.10.10.10","10.10.10.11","127.0.0.1"],
    "HostPath" : "/home/shell/shell.sh",
    "HotsUser" : "shell",
    "CrontabEnable" : true,
    "CrontabData" : {
        "Time" : "0 1 * * *",
        "command" : "/usr/bin/sh",
        "arg" : "/home/shell/shell.sh"
    },
    "Language" : "bash",
    "Authorer" : "luojiashuo",
    "Description" : "shell脚本注释模板"
}
COMMENT_BLOCK

# 测试脚本
date=$(date "+%FT%T.000")
echo $date
