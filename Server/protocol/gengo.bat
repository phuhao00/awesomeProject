@echo off
echo ��������Go����...
 protoc -I=./ --go_out=../PB/  .\messageid.proto
pause