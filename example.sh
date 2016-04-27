#!/bin/bash
set -x
go run echo_server.go -bind :8001 &
go run echo_server.go -bind :8002 &
go run echo_server.go -bind :8003 &
go run echo_server.go -bind :8004 &
go run tcptee.go -bind :8000 -backends :8001,:8002,:8003,:8004 &
sleep 2
echo 'All backend echo_servers should echo this message!' | nc -w0 127.0.0.1 8000
kill $(lsof -i:8000 -i:8001 -i:8002 -i:8003 -i:8004 -t)
sleep 2
jobs
