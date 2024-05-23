: <<COMMENT_BLOCK
{
    "HostLocation" : "HuaweiCloud",
    "HostIP" : ["10.10.10.10","10.10.9.70","127.0.0.1"],
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
COMMENT_BLOCK

# 测试脚本
date=$(date "+%FT%T.000")
echo $date
echo "this is shell test"
