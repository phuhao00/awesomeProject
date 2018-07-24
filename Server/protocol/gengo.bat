@echo off
echo 正在生成Go代码...
 protoc -I=./ --go_out=../PB/  .\messageid.proto
pause